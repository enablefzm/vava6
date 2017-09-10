package vatools

import (
	"fmt"
	"sync"
)

type VaMap struct {
	mp map[interface{}]interface{}
	lk *sync.RWMutex
	fc func(key interface{}) (interface{}, error)
}

func NewVaMap(max int) *VaMap {
	return &VaMap{
		mp: make(map[interface{}]interface{}, max),
		lk: new(sync.RWMutex),
	}
}

func (this *VaMap) Put(key, value interface{}) {
	this.lk.Lock()
	this.mp[key] = value
	this.lk.Unlock()
}

func (this *VaMap) Get(key interface{}) (value interface{}, ok bool) {
	this.lk.RLock()
	value, ok = this.mp[key]
	this.lk.RUnlock()
	return
}

// 通过缓存获取指定Key的对象
//	@parames
//		key			interface{}									要获取的关键key
//		fnCreate	func(k interface{}) (interface{}, error)	如果没有这个对象则创建这个对象的方法
//	@return
//		interface{}
//		error
func (this *VaMap) GetCache(key interface{}, fnCreate func(k interface{}) (interface{}, error)) (interface{}, error) {
	if res, ok := this.Get(key); ok {
		return res, nil
	}
	// 创建对象
	this.lk.Lock()
	// 在锁住状态进行第二次获取
	res, ok := this.mp[key]
	if ok {
		this.lk.Unlock()
		return res, nil
	}
	// 第二次还是没有获取到则创建一个新对象
	v, err := fnCreate(key)
	if err == nil {
		this.mp[key] = v
	}
	this.lk.Unlock()
	return v, err
}

func (this *VaMap) GetCacheOnFunc(key interface{}) (interface{}, error) {
	if this.fc == nil {
		return nil, fmt.Errorf("No func")
	}
	return this.GetCache(key, this.fc)
}

func (this *VaMap) SetFc(fc func(key interface{}) (interface{}, error)) {
	this.fc = fc
}

func (this *VaMap) Len() int {
	return len(this.mp)
}
