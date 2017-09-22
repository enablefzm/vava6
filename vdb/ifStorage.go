package vdb

// 执行保存的对象
type IFStorage interface {
	// 查询数据
	//	@parames
	//		talbe 	string		要被查询的表
	//		keys	[]string	要被查询的关键KEY
	//	@result
	//		[]map[string]string
	//		error
	Querys(table string, keys map[string]interface{}) ([]map[string]string, error)

	// 插入数据
	//	@parames
	//		table 	插入保存到哪个表
	//		info	要被插入的数据
	//	@result
	//		int64	返回新增的ID
	//		error	错误信息
	Insert(table string, info map[string]interface{}) (int64, error)

	// 更新记录
	//	@parames
	//		table	更新的数据表
	//		info	要被保存的数据
	//		keys	更新的条件
	//	@return
	//		int64	受影响的行数
	//		error	错误对象
	Update(table string, info, keys map[string]interface{}) (int64, error)

	// 删除记录
	//	@parames
	//		table	删除的数据表
	//		keys	删除的条件
	//	@return
	//		int64	受影响的行数
	//		error	错误对象
	Remove(table string, keys map[string]interface{}) (int64, error)
}
