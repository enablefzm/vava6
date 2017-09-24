package vdb

func NewSaveInfo(tlbName string, mpKeys, mpSave map[string]interface{}, isNew bool) *SaveInfo {
	return &SaveInfo{
		TableName: tlbName,
		MpKeys:    mpKeys,
		MpSave:    mpSave,
		IsNew:     isNew,
	}
}

type SaveInfo struct {
	TableName string
	MpKeys    map[string]interface{}
	MpSave    map[string]interface{}
	IsNew     bool
}

// 被保存的对象信息
type IFTargetDB interface {
	// 获取要被存储的信息
	GetSaveDB() *SaveInfo

	// 设定为不是新增对象
	SetIdx(id uint)
}
