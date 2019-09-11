package dao

import (
	"fmt"
	"database/sql"   //sql抽象层包
	_"github.com/go-sql-driver/mysql"  //mysql驱动包(匿名导入，只为使用mysql标识符)
	"strings"
	"strconv"
)



//定义原生字符串类型，不对函数做反引号处理（count(),sum()...)
type Raw struct{
	String string
}

// 定义结构体对象
type Dao struct{
	db *sql.DB
	sql string
	sqlParams []interface{}
	// query的每个部分
	queryTable string	//表名
	queryCondition string	//where条件
	queryParams []interface{}	//where条件值部分
	queryDistinct bool	//去重
	queryFields string	//查询的字段列表部分
	queryJoins []string //join A... join B...
	queryGroupBy string
	queryHavingCondition string	//字段条件部分 id<?
	queryHavingParams []interface{}	//条件？值部分
	queryOrderBy string	//字段排序方式部分
	queryLimit string	//限定结果集数量部分
	queryOffset string	//偏移量部分
}

//级联调用，表名
func (this *Dao) Table(table string) *Dao {
	//调用私有表名方法
	this.queryTable = _table(table)
	//返回当前对象
	return this
}



//from方法与表名一致
func (this *Dao) From(table string) *Dao {
	return this.Table(table)
}

//级联调用，条件, condition（代表条件部分，例：id> ）,params（代表条件参数值部分 ）
func (this *Dao) Where(condition string, params []interface{}) *Dao{
	this.queryCondition, this.queryParams = condition, params
	return this
}



// join表连接子句部分, select * from A表 as 别名 join B表 as 别名 on 连接条件;
func (this *Dao) Join(joinType string, joinTable string, joinOn string) *Dao{
	//将传递的 连接类型 统一转成大写
	join := fmt.Sprintf("%s JOIN %s ON %s", strings.ToUpper(joinType), _table(joinTable), joinOn)
	//将拼凑好的子句，追加写入切片中
	this.queryJoins = append(this.queryJoins, join)
	return this
}

//左连接
func (this *Dao) LeftJoin(joinTable string, joinOn string) *Dao{
	return this.Join("LEFT", joinTable, joinOn)
}
//右连接
func (this *Dao) RightJoin(joinTable string, joinOn string) *Dao{
	return this.Join("RIGHT", joinTable, joinOn)
}
//内连接
func (this *Dao) InnerJoin(joinTable string, joinOn string) *Dao{
	return this.Join("INNER", joinTable, joinOn)
}

//去重复 distinct *, (select distinct * from test1;)
func (this *Dao) Distinct() *Dao{
	this.queryDistinct = true
	return this
}

//GroupBy 分组子句  group by r.rank , br.rank
func (this *Dao) GroupBy(fields ...string) *Dao{
	fieldSlice := []string{}
	//可能存在多字段，需要遍历获取每个参数
	for _, field := range fields{
		//反引号包裹标识符
		fieldSlice = append(fieldSlice,_fieldwrap(field))
	}
	//多字段间用逗号分隔
	this.queryGroupBy = strings.Join(fieldSlice, ",")
	return this
}

//Having 结果筛选 子句
func (this *Dao) Having(condition string, params []interface{}) *Dao{
	this.queryHavingCondition, this.queryHavingParams = condition, params
	return this
}

//OrderBy() 字段排序子句
func (this *Dao) OrderBy(fields ...string) *Dao{
	fieldSlice := []string{}
	//考虑字段别名问题，先遍历,然后为表名和别名使用反引号包裹
	for _, field := range fields{
		//判断是否存在空格，存在则说明指点了排序方式desc，asc
		if fieldSplit := strings.Split(field, " "); len(fieldSplit) > 1{
			// order by `zj_id` DESC || ASC;
			fieldSlice = append(fieldSlice, _fieldwrap(fieldSplit[0]) +" "+ strings.ToUpper(fieldSplit[1]))
		}else{
			//不存在空格，说明没有指定排序方式
			fieldSlice = append(fieldSlice, _fieldwrap(fieldSplit[0]))
		}
	}
	//若存切片内存在多个值，用逗号分隔
	this.queryOrderBy = strings.Join(fieldSlice, ",")
	return this
}

// Limit() 限制结果集数量 子句
func (this *Dao) Limit(size int) *Dao{
	// strconv.Itoa() 将int整形转为字符串类型
	this.queryLimit = strconv.Itoa(size)
	return this
}
//Offset 偏移量，配合limit使用
func (this *Dao) Offset(rows int) *Dao{
	this.queryOffset = strconv.Itoa(rows)
	return this
}

// fields字段列表，字段列表多样性: 字段应用反引号
// （* , `id`, `name` as `na`, `name` `na`, `a`.`password`, sum(`field`), `b`.`class` as `b_class` ) 
//调用时传递， this.Fields("*", "name", "name as na", "name na", "a.title")
func (this *Dao) Fields(fields ...interface{}) *Dao{
	fieldList := []string{}
	//遍历参数，获得每个字段
	for _, field := range fields {
		//判断类型
		switch field.(type) {
		case string:					
			//考虑每个字段的情况
			//是否包含空格，有空格意味着有别名,断言字符串类型
			fieldSlice := strings.Split(field.(string), " ") //如果这个参数有空格，则拆分。得到切片格式： []string{name as na}
			//判断这个切片长度，如果 >1 说明有空格是别名，<1 则无空格
			if fslen := len(fieldSlice); fslen> 1 {
				//存在空格，存在别名，全部处理成带有 as 格式
				fieldList = append(fieldList,fmt.Sprintf("%s AS `%s`", _fieldwrap(fieldSlice[0]), fieldSlice[fslen-1]))	//存在空格的切片中，首个参数为字段，最后为别名
			}else{
				//无别名的情况
				fieldList = append(fieldList, _fieldwrap(fieldSlice[0]))
			}
		//类型为Raw时，说明是函数，不进行加反引号处理
		case Raw:
			fieldList = append(fieldList, field.(Raw).String)	//断言为Raw原生字符串类型
	}
}
	//把处理好的有反引号包裹的字段，使用逗号分隔（`f1` AS `a1`,*,`t`.`wq` AS `t1`,`t`.*）
	this.queryFields = strings.Join(fieldList, ",")
	return this
}


//Column()	查询第一条记录的第一个字段的值返回
func (this *Dao) Column() ([]string, error){
	rows, cols, err := this.rowsAndCols()
	if err != nil{
		return []string{}, err
	}
	//获取第一列的数据
	firstCol := cols[0]
	result := []string{}
	for _, row := range rows{
		result = append(result, row[firstCol])
	}
	return result, nil
}

//fetchRow 查询单行
func (this *Dao) Row() (map[string]string, error){
	result, _, err := this.rowAndCols()
	return result, err
}


//查询单个字段值
func (this *Dao) Value() (string, error){
	row, cols, err := this.rowAndCols()
	if err != nil{
		return "", err
	}
	firstCol := cols[0]
	if len(row) > 0{
		//返回第一个列
		return row[firstCol],nil
	}else{
		return "", nil
	}
	return "", nil
}





//fetchAll查询全部（多行）
func (this *Dao) Rows() ([]map[string]string, error){
	result, _, err := this.rowsAndCols()
	// if err != nil{
	// 	return []map[string]string{}, err
	// }
	return result, err
}

//实现insert插入方法
func (this *Dao) Insert(fields map[string]interface{})(int64,error){
	//准备工作，遍历fields，得到全部key（字段名）和value(值)
	fieldList, valueList, valuePLList := []string{}, []interface{}{},[]string{}
	for field, value := range fields{
		//将遍历出来的key和value分别存储，字段名需要使用反引号包裹
		fieldList = append(fieldList, "`"+ field +"`")
		valueList = append(valueList, value)
		//值的占位符(遍历多少次值，就得到多少个？)
		valuePLList = append(valuePLList,"?")
	}
	//拼凑SQL语句，insert into 表名(字段列表) values(值列表)
	// query := "INSERT INTO"  //固定部分
	// //表名
	// query += "`" + table + "`"     //==insert into `test1`
	// //字段 (`id`,`username`,`possword`)
	// query += " "+ "("+ strings.Join(fieldList, ",") +")"
	// // 值 values(?,?)  , 使用？占位
	// query += " "+ "VALUES"+ " "+ "("+ strings.Join(valuePLList,",") +")"

	// 以上拼凑可以使用fmt.Sprintf()函数实现，格式化字符串返回
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		this.queryTable,
		strings.Join(fieldList, ","),
		strings.Join(valuePLList,","))

		// 执行,将值列表展开传递给？
	// result, err := this.db.Exec(query, valueList...)
	//调用执行非查询类方法，记录SQL语句
	result, err := this.exec(query, valueList...)
	//执行过后，清理掉记录，防止下次调用冲突
	this.clearQuery()
	if err != nil{
		return 0,err
	}
	// 执行成功，返回lastInsertID
	id, err := result.LastInsertId()
	if err != nil{
		//没有拿到id也不是错误,考虑可能没有id自增的字段存在
		return 0,nil		
	}
	//返回
	return id,nil
}

//实现update更新方法
// update 表名 set 字段 = 新值 where id=1;
func (this *Dao) Update(fields map[string]interface{})(int64,error){
	// 准备工作
	// set `field` = ?
	setList, valueList := []string{}, []interface{}{}
	for field, value := range fields{
		//条件字段列表
		setList = append(setList, fmt.Sprintf("`%s`=?",field))
		//条件值列表
		valueList = append(valueList,value)
	}
	//拼凑
	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s",this.queryTable, strings.Join(setList,","), this.queryCondition)
	//执行SQL , 向？传递的数据部分由 字段和条件参数组成
	// result, err := this.db.Exec(query,append(valueList, this.queryParams...)...) 
	//调用执行非查询类方法，记录SQL语句
	result, err := this.exec(query,append(valueList, this.queryParams...)...) 
	//执行过后，清理掉记录，防止下次调用冲突
	this.clearQuery()
		if err != nil{
			return 0,err
		}
	//执行成功,返回受影响数量
	rows, err := result.RowsAffected()
	if err != nil{
		return 0,err
	}
	return rows, nil
}

//Delete删除， delete from 表名 where condition
func (this *Dao) Delete()(int64, error){
	//query := delete from test1 from where id>? 
	query := fmt.Sprintf("DELETE FROM %s WHERE %s", this.queryTable, this.queryCondition)
	//执行SQL
	// result, err := this.db.Exec(query, this.queryParams...)
	//调用执行非查询类方法，记录SQL语句
	result, err := this.exec(query, this.queryParams...)
	//执行过后，清理掉记录，防止下次调用冲突
	this.clearQuery()
	if err != nil{
		return 0,err
	}
	//执行成功,返回受影响数量
	rows, err := result.RowsAffected()
	if err != nil{
		return 0,err
	}
	return rows, nil
}

//获取SQL语句和参数
func (this *Dao) LastSQL() (string, []interface{}){
	// 一次性获取
	sql, params := this.sql, this.sqlParams
	this.sql, this.sqlParams = "", []interface{}{}
	return sql, params
}

//构造函数, 初始化数据库，连接，测试
//参数为字符串映射表， [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
func NewDao(config map[string]string)(*Dao, error){
	//一，定义变量来存储最终拼凑DSN成的SQL语句
		DSN := ""
		// 初始化服务器的连接,判断是否存在相应的元素
		username, ok := config["username"]
		//考虑用户没有传递服务器用户名的情况
		if !ok{
			// 若未传递用户名，则设置默认值
			username = ""
		}
	
		password, ok := config["password"]
		if !ok{
			password = ""
		}
	
		//拼凑[username[:password]@]
		if "" == username{
			// 判断若用户名为空，则将字符串加到DSN当中
			DSN += ""
		}else if "" == password{
			//若用户名存在，密码为空
			DSN += username + "@"
		}else{
			//若都存在
			DSN += username + ":" + password + "@"
		}
	
		//拼凑[protocol[(address)]]
		protocol, ok := config["protocol"]
		if !ok {
			// 若为空，则设置默认值tcp网络连接协议
			protocol = "tcp"
		}
		host, ok := config["host"]
		if !ok {
			host = "localhost"
		}
		//若为空，则设置默认端口3306
		port, ok := config["port"]
		if !ok {
			port = "3306"
		}
		DSN += protocol + "(" + host + ":" + port + ")"
	
		// 判断是否有传递数据库名，/dbname
		dbname, ok := config["dbname"]
		if !ok {
			// 若为空，则设置默认值为空
			dbname = ""
		}
		DSN += "/" + dbname
	
		//判断是否有传递字符集，校对集，[?param1=value1&...&paramN=valueN]
		//param可以有多条组成，所以使用变量存储
		params := []string{}
		collation, ok := config["collation"]
		if !ok {
			// 若为空，则设置默认值为空
			collation = ""
		}else{
			//若有传递值，则将值追加到切片中
			params = append(params, "collation="+ collation)
		}
		//使用.Join将切片中的多个值，用&符号连接
		paramStr := strings.Join(params, "&")
	
		//最终拼凑
		DSN += "?" + paramStr
	
	// 二， 连接mysq
		db, err := sql.Open("mysql",DSN)
		if err != nil {
			//若连接错误返回空对象，和错误提示
			return nil,err
		}
		// 测试是否连接成功
		if pingErr:= db.Ping(); pingErr != nil{
			return nil,pingErr
			
		}
	
	// 三，返回Dao对象
		// 实例化对象
		this := new(Dao)
		this.db = db
		//初始化query部分
		this.clearQuery()
		return this, nil
	}


