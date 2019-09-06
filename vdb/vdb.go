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
	saveInfo := ptrTarget.GetSaveDB()
	if saveInfo.IsNew {
		if idx, err := this.Insert(saveInfo.TableName, saveInfo.MpSave); err != nil {
			return err
		} else {
			ptrTarget.SetIdx(uint(idx))
		}
	} else {
		for k, _ := range saveInfo.MpKeys {
			delete(saveInfo.MpSave, k)
		}
		if rows, err := this.Update(saveInfo.TableName, saveInfo.MpSave, saveInfo.MpKeys); err != nil {
			return err
		} else if rows < 1 {
			// 更新不成功则尝试新增
			for k, v := range saveInfo.MpKeys {
				saveInfo.MpSave[k] = v
			}
			if _, err := this.Insert(saveInfo.TableName, saveInfo.MpSave); err != nil {
				return err
			}
		}
	}
	return nil
}
