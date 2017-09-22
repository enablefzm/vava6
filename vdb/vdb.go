package vdb

type DBStorage struct {
	IFStorage
}

func NewDBStorage(ptr IFStorage) *DBStorage {
	return &DBStorage{
		IFStorage: ptr,
	}
}

// 保存对象
//	@parames
//		ptrTarget 	要被保存的对象
//	@return
//		error
func (this *DBStorage) Save(ptrTarget IFTargetDB) error {
	if ptrTarget.IsNew() {
		if idx, err := this.Insert(ptrTarget.GetTableName(), ptrTarget.GetSave()); err != nil {
			return err
		} else {
			ptrTarget.SetIdx(int(idx))
		}
	} else {
		if _, err := this.Update(ptrTarget.GetTableName(), ptrTarget.GetSave(), ptrTarget.GetKeys()); err != nil {
			return err
		}
	}
	return nil
}
