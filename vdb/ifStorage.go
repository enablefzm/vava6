package vdb

// 被保存的对象信息
type IFTargetDB interface {
	GetSave() map[string]interface{}
}
