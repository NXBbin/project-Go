
export default {
    // 存储定义数据
    state: {
        products: [
            // 示例数据
            // {productID: 1, buyQuantity: 12,},
            // {productID: 11, BuyQuantity: 1,}
        ],
    },
    // 获取数据
    getters: {
        products: state => state.products
    },
    // 更新数据
    mutations: {
        addProduct(state, buyInfo) {
            state.products.push(buyInfo)
        },
        addProductBuyQuantity(state, product) {
            state.products[product.index].buyQuantity += product.buyQuantity
        },
        setProducts(state, payload) {
            state.products = payload.products
        },
        setProduct(state, payload) {
            state.products[payload.index] = payload.buyInfo
        },
        removeProduct(state, payload) {
            state.products.splice(payload.index, 1)
        }
    },
    // 触发更新操作
    actions: {
        // 添加商品到购物车
        addToCart(context, buyInfo) {
            // 判断产品是否存在与购物车中
            let index = context.getters.products.findIndex(element=>{
                return element.productID == buyInfo.productID
            })

            if (-1 == index) {
                // 该产品不存在，应该添加
                context.commit('addProduct', buyInfo)
            }
            else {
                // 产品存在
                context.commit('addProductBuyQuantity', {index, buyQuantity: buyInfo.buyQuantity})
            }

            // 持久化存储
            context.dispatch('save')
        }, 

        // 将产品从购物车中移除
        removeFromCart(context, payload) {
            // 找到需要删除的索引
            let index = context.getters.products.findIndex(element=>{
                return element.productID == payload.productID
            })
            // 利用mutation完成更新
            context.commit('removeProduct', {
                index,
            })

            // 持久化存储
            context.dispatch('save')
        },

        // 设置购买数量
        setBuyQuantity(context, payload) {
            let index = context.getters.products.findIndex(element=>{
                return element.productID == payload.productID
            })
             // 利用mutation完成更新
             context.commit('setProduct', {
                index,
                buyInfo: payload,
            })
            // 持久化存储
            context.dispatch('save')
        },

        // 统计购物车信息
        cartInfo(context, payload) {
            // 遍历全部产品
            let buyQuantityTotal = 0
            for (let p of context.getters.products) {
                buyQuantityTotal += p.buyQuantity
            }

            return {
                buyQuantityTotal,
            }
        },

        // 持久化存储
        save(context, payload) {
            // 在未登录的情况下
            window.localStorage.setItem('cart', JSON.stringify(context.getters.products))
        },
        // 初始化购物车
        init(context, payload) {
            let products = JSON.parse(window.localStorage.getItem('cart')) || []
            context.commit('setProducts', {products})
        },
    },
  }