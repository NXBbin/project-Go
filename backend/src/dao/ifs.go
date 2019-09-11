package dao



//定义该包所要实现的接口方法
type I_Dao interface{	
	Distinct() *Dao
	Fields(fields ...interface{}) *Dao
	Table( string) *Dao
	From( string) *Dao
	// join表连接方法 (连接类型[左右内连接]，连接表名，连接条件)
	Join( string,  string,  string) *Dao
	LeftJoin( string,  string) *Dao
	RightJoin( string,  string) *Dao
	InnerJoin( string,  string) *Dao
	Where( string,  []interface{}) *Dao
	GroupBy( ...string) *Dao
	Having( string,  []interface{}) *Dao
	OrderBy( ...string) *Dao
	Limit( int) *Dao
	Offset( int) *Dao

	Column() ([]string, error)
	Row() (map[string]string, error)
	Value() (string, error)
	// 多行查询
	Rows() ([]map[string]string, error)
	//SQL插入方法,参数：表名和值列表,返回：最新生成的ID和错误
	Insert(  map[string]interface{})(int64,error)	
	//update更新方法, 注意where条件部分的参数，condition代表:username like ? and id? , params代表：[]{"bin", "1"}
	Update(  map[string]interface{})(int64,error)	
	Delete()(int64, error)

	LastSQL() (string, []interface{})
}