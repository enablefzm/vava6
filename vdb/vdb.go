package vdb

// 被保存的对象信息
type IFTargetDB interface {
	GetSave() map[string]interface{}
}

// 执行保存的对象
type IFStorage interface {
	// 查询数据
	Querys(table string, key []string) ([]map[string]string, error)
}
