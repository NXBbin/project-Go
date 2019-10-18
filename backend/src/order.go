package main

import (
	"config"
	"controller"
	"encoding/json"
	"log"
	"model"
	"time"
)

// 订单处理进程

func main() {
	log.Println("订单处理中...")

	//初始化配置
	config.InitConfig()

	// 调用初始化MYSQL，（连接数据库，gorm）
	db, err := controller.InitDB()
	if err != nil {
		log.Println("数据库连接失败", err)
		return
	}
	defer db.Close()

	//初始化redis
	rc, err := controller.InitRedis()
	if err != nil {
		log.Println(err)
		return
	}
	defer rc.Close()

	//从消息队列中，获取订单信息，进行处理
	for {
		//循环3秒读取最新一条库存信息
		result, err := controller.Rds.Do("XREAD", "COUNT", "1", "BLOCK", "3000", "STREAMS", "orderQueue", "$")
		if err != nil {
			log.Println(err)
			continue
		}
		if result == nil {
			//当前没有新订单，继续等待
			continue
		}
		// 解析result数据，一级级的断言获取订单号sn
		sn := string(result.([]interface{})[0].([]interface{})[1].([]interface{})[0].([]interface{})[1].([]interface{})[1].([]byte))
		log.Println(sn)

		// [
		// 	[
		// 		[111 114 100 101 114 81 117 101 117 101]
		// 		[
		// 			[
		// 				[49 53 55 49 51 51 56 56 51 48 52 52 53 45 48]
		// 				[
		// 					[99111 110 116 101 110 116]
		// 					 [50 48 49 57 49 48 49 55 48 48 48 48 48 48 48 56]
		// 				]
		// 			]
		// 		]
		// 	]
		// ]

		//获取订单信息
		orderResult, err := controller.Rds.Do("HGET", "tempOrder", sn)
		if err != nil {
			log.Println(err)
			continue
		}
		orderInfo := controller.TempOrder{}
		parseErr := json.Unmarshal(orderResult.([]byte), &orderInfo)

		if parseErr != nil {
			log.Println(parseErr)
			continue
		}
		log.Println(orderInfo)

		//获取用户购物车中的全部购买商品和数量
		cart := model.Cart{}
		db.Where("user_id=?", orderInfo.UserID).Find(&cart)
		cartInfo := []struct {
			ProductID   uint `json:"productID"`
			BuyQuantity int  `json:"buyQuantity`
		}{}
		//解析JSON数据，
		json.Unmarshal([]byte(cart.Content), &cartInfo)
		// log.Println("cartInfo数据:", cartInfo)
		// log.Println("cart.Content数据:", cart.Content)

		//核实订单信息，核实购买产品和购买数量
		buyInfo := []struct {
			ProductID   uint `json:"productID"`
			BuyQuantity int  `json:"buyQuantity`
		}{}
		// 遍历得到购物车的产品ID和数量
		for _, p := range cartInfo {
			// log.Println("p数据（购物车中产品ID和数量）:", p)
			// 遍历得到订单购买的产品ID
			for _, pID := range orderInfo.BuyProductID {
				// log.Println("pID数据:", pID)
				//如果购物车中产品ID和订单购买产品ID相同时，则将购买产品和数量加入队列
				if p.ProductID == pID {
					buyInfo = append(buyInfo, p)
					break
				}
			}
		}
		//检测所需要的购买的产品是否存在于购物车中
		if len(buyInfo) != len(orderInfo.BuyProductID) {
			log.Println("所购买的产品不再购物车中")
			controller.Rds.Do("HSET", "orderResult", sn, "error")
			continue
		}

		//检查库存
		flag := true
		for _, p := range buyInfo {
			//获取每个产品对应的库存，进行检测
			quantity := model.Quantity{}
			db.Where("product_id=?", p.ProductID).Find(&quantity)
			if quantity.Number < p.BuyQuantity {
				//系统库存，小于购买库存，则失败
				flag = false
				break
			}
		}
		if !flag {
			//库存不足，订单失败，记录结果
			log.Println("订单失败，库存不足")
			controller.Rds.Do("HSET", "orderResult", sn, "error")
			continue
		}
		//库存足够
		log.Println("订单成功，库存充足")

		//扣减库存
		for _, p := range buyInfo {
			//获取每个产品对应的库存，进行检测
			quantity := model.Quantity{}
			db.Where("product_id=?", p.ProductID).Find(&quantity)
			//扣减库存（更新操作）总库存-购买数量
			quantity.Number = quantity.Number - p.BuyQuantity
			db.Save(&quantity)
		}

		//形成订单数据，插入订单表
		order := model.Order{}
		order.Sn = sn
		order.AddressID = orderInfo.AddressID
		order.PaymentStatusID = 2 //未支付
		order.ShippingID = orderInfo.ShippingID
		order.ShippingStatusID = 2 //未发货
		order.Amount = 99          //总金额
		order.OrderTime = time.Now()
		order.OrderStatusID = 2 //确认
		// order.UserID = orderInfo.UserID
		db.Create(&order)

		//记录订单中的产品
		for _, p := range buyInfo {
			//获取订单中的产品信息
			product := model.Product{}
			db.Where("id", p.ProductID).Find(&product)

			//形成订单数据，插入订单产品关联表
			orderProduct := model.OrderProduct{}
			orderProduct.OrderID = order.ID
			orderProduct.ProductID = p.ProductID
			orderProduct.BuyQuantity = p.BuyQuantity
			orderProduct.BuyPice = int(product.Price * 100)
			db.Create(&orderProduct)
		}

		//删除购物车已购的产品
		restInfo := []struct {
			ProductID   uint `json:"productID"`
			BuyQuantity int  `json:"buyQuantity`
		}{}
		// 遍历得到购物车的产品ID和数量
		for _, p := range cartInfo {
			// log.Println("p数据（购物车中产品ID和数量）:", p)
			// 遍历得到订单购买的产品ID
			flag := true //假定为需要保留
			for _, pID := range orderInfo.BuyProductID {
				// log.Println("pID数据:", pID)
				//如果购物车中产品ID和订单购买产品ID相同时
				if pID == p.ProductID {
					flag = false //该产品已经购买的，需要从购物车中删除
					break
				}
			}
			if flag {
				restInfo = append(restInfo, p)
			}
		}
		//更新购物车字段
		c, _ := json.Marshal(restInfo)
		cart.Content = string(c)
		db.Save(&cart)

		//记录成功结果
		controller.Rds.Do("HSET", "orderResult", sn, "success")
	}
}
