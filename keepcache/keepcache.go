package keepcache

import (
	"sync"
)

func NewKeepCache(fnCreate func(k interface{}) (interface{}, error)) *KeepCache {
	return &KeepCache{
		mpCache:    make(map[interface{}]interface{}, 100),
		lk:         new(sync.RWMutex),
		createFunc: fnCreate,
	}
}

type KeepCache struct {
	mpCache    map[interface{}]interface{}
	lk         *sync.RWMutex
	createFunc func(k interface{}) (interface{}, error)
}

func (this *KeepCache) Get(k interface{}) (interface{}, error) {
	var err error
	this.lk.RLock()
	v, ok := this.mpCache[k]
	this.lk.RUnlock()
	if !ok {
		this.lk.Lock()
		if v, ok = this.mpCache[k]; !ok {
			v, err = this.createFunc(k)
			if err == nil {
				this.mpCache[k] = v
			}
		}
		this.lk.Unlock()
	}
	return v, err
}
