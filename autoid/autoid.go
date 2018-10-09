package autoid

import (
	"sync"
	"vava6/vatools"
)

func NewAutoID() *AutoID {
	return &AutoID{
		cacheID: make([]int, 0, 20),
		lk:      new(sync.Mutex),
	}
}

func NewAutoIDOnJson(strJson string) *AutoID {
	db := &AutoIdDB{}
	vatools.UnJson(strJson, db)
	return NewAutoIDOnDb(db)
}

func NewAutoIDOnDb(db *AutoIdDB) *AutoID {
	return &AutoID{
		cacheID: db.CacheID,
		lastID:  db.LastID,
		lk:      new(sync.Mutex),
	}
}

type AutoID struct {
	lastID  int
	cacheID []int
	lk      *sync.Mutex
}

// 获取ID
func (this *AutoID) GetID() int {
	this.lk.Lock()
	defer this.lk.Unlock()
	// 判断缓存里是否有可用ID
	if len(this.cacheID) > 0 {
		id := this.cacheID[0]
		this.cacheID = this.cacheID[1:]
		return id
	}
	this.lastID++
	return this.lastID
}

// 将ID放回缓存，供下次复用
func (this *AutoID) PutID(id int) {
	this.lk.Lock()
	this.cacheID = append(this.cacheID, id)
	this.lk.Unlock()
}

func (this *AutoID) GetJsonSave() string {
	db := this.GetDbSave()
	strJson, _ := vatools.Json(db)
	return strJson
}

func (this *AutoID) GetDbSave() *AutoIdDB {
	return &AutoIdDB{
		CacheID: this.cacheID,
		LastID:  this.lastID,
	}
}
