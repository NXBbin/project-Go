package dao

import (
	"fmt"
	"database/sql"   //sql抽象层包
	"strings"

)



//单行和列
func (this *Dao) rowAndCols() (map[string]string, []string, error){
	//调用limit，保证只查询到一条数据
	this.Limit(1)
	rows, cols, err := this.rowsAndCols()
	if err != nil{
		return map[string]string{}, []string{} , err
	}

	//判断是否查询到数据
	if len(rows) > 0 {
		return rows[0], cols, nil
	}else{
		//一条数据都不存在,但不是错误
		return map[string]string{}, []string{} ,  nil
	}
}

//buildSelect  拼凑SQL语句
func (this *Dao) buildSelect() string{
	query := "SELECT"
	//判断是否有distinct去重子句出现
	if this.queryDistinct{
		query += " DISTINCT"
	}
	//字段部分
	query += " "+ this.queryFields
	//表名部分
	query += " FROM "+ this.queryTable
	// Join部分,切片类型，可能存在多个值，需要遍历得到
	for _, join := range this.queryJoins{
		query += " "+ join
	}
	// Where部分
	if "" != this.queryCondition{
		query += " WHERE " + this.queryCondition
	}
	// groupby部分
	if "" != this.queryGroupBy{
		query += " GROUP BY " + this.queryGroupBy
	}
	//having 部分
	if "" != this.queryHavingCondition{
		query += " HAVING " + this.queryHavingCondition
	}
	//order by 部分
	if "" != this.queryOrderBy{
		query += " ORDER BY " + this.queryOrderBy
	}
	// limit 和 offset 部分
	//判断是否同时存在limit 和 offset 子句
	if "" != this.queryOffset && "" != this.queryLimit{
		query += " LIMIT " + this.queryLimit + " OFFSET " + this.queryOffset
	}else if "" != this.queryLimit{
		query += " LIMIT " + this.queryLimit
	}
	// fmt.Println(query)
	return query
}

//定义rowsAndCols 用于获取列的数据
func (this *Dao) rowsAndCols() ([]map[string]string, []string, error){
	//构建SQL
	query := this.buildSelect()
	// 执行
	// rows, err := this.db.Query(query, append(this.queryParams,this.queryHavingParams...)...)
	//调用执行非查询类方法，记录SQL语句
	rows, err := this.query(query, append(this.queryParams,this.queryHavingParams...)...)
	this.clearQuery()
	if err != nil{
		return []map[string]string{}, []string{}, err
	}
	defer rows.Close()
	//确定列的数量(字段数量)
	cols, colError := rows.Columns()
		if colError != nil{
			return []map[string]string{}, []string{}, colError
		}
	//确定字段数量
	colNum := len(cols)
	//为了rows.Scan传参，不定数量的接口切片类型，并且需要展开
	fields := make([]interface{}, colNum)
	//values存储的是字符串的引用
	values := make([]sql.NullString, colNum)
	for i,_ := range fields{
		//保证fields 的每个元素都是指针类型，可以被赋值*string
		fields[i] = &values[i]
	}

	result := []map[string]string{}
	//处理结果
	for rows.Next(){
		scanErr := rows.Scan(fields...)
		if scanErr != nil{
			//当前记录scan错误（可能出现NULL），继续获取下一条记录
			continue 
		}
		//获取到数据，整理成目标格式
		row := map[string]string{}
		for i,_ := range fields{
			//通过接口断言，再解析地址的方式得到字符串
			// row[cols[i]] = *(fields[i].(*string))
			ns :=  *(fields[i].(*sql.NullString))
			//判断是否为合法数据，为NULL则输出字符串“//null”
			if ns.Valid{
				row[cols[i]] = ns.String
			}else{
				row[cols[i]] = "//null"
			}
		}
		result = append(result, row)
	}
	return result, cols, nil
}

// 清空记录query
func (this *Dao) clearQuery(){
	this.queryTable = ""
	this.queryCondition = ""
	this.queryParams = []interface{}{}
	this.queryDistinct = false
	this.queryFields = "*"
	this.queryJoins = []string{}
	this.queryGroupBy = ""
	this.queryHavingCondition = ""
	this.queryHavingParams = []interface{}{}
	this.queryOrderBy = ""
	this.queryLimit = ""
	this.queryOffset = ""
}

//内部表名函数，用于其他子句方法调用表名
func _table(table string) string{
	//是否包含空格，有空格意味着有别名
	tableSlice := strings.Split(table, " ")
	if tslen := len(tableSlice); tslen> 1 {
		// 存在空格,统一处理成 as 别名格式
		return fmt.Sprintf("`%s` AS `%s`",tableSlice[0], tableSlice[tslen-1])
	}else{
		//不存在空格，说明只有一个表名，直接反引号包裹即可
		return fmt.Sprintf("`%s`",tableSlice[0])
	}
}

//内部函数，判断字段是否有.部分,存在.就分别包裹反引号
func _fieldwrap(field string) string{
	fieldSlice := strings.Split(field, ".")
	if fslen := len(fieldSlice); fslen> 1 {
		// 判断多表*查询 t.*
		if fieldSlice[1] == "*" {
			return fmt.Sprintf("`%s`.%s",fieldSlice[0], fieldSlice[1])
		}else{
			return fmt.Sprintf("`%s`.`%s`",fieldSlice[0], fieldSlice[1])
		}
	}else{
		//判断字段起始就为 * 的情况
		if fieldSlice[0] == "*"{
			return fmt.Sprintf("%s",fieldSlice[0])
			//判断是否为函数(sum(),contains()...)
		}else{
			return fmt.Sprintf("`%s`",fieldSlice[0])
		}
	}
}


//记录SQL语句,执行非查询类
func (this *Dao) exec(query string, params ...interface{}) (sql.Result, error){
	this.sql = query
	this.sqlParams = params
	return this.db.Exec(query, params...)
}
//执行查询类
func (this *Dao) query(query string, params ...interface{}) (*sql.Rows, error){
	this.sql = query
	this.sqlParams = params
	return this.db.Query(query, params...)
}











