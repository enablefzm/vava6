package vdb

// 被保存的对象信息
type IFTargetDB interface {
	// 获取要被存储的信息
	GetSave() map[string]interface{}

	// 获取TableName
	GetTableName() string

	// 获取索引Key
	GetKeys() map[string]interface{}

	// 判断是不是新增对象
	IsNew() bool

	// 设定为不是新增对象
	SetIdx(id int)
}
