package pools

import (
	"fmt"
	"sync"
)

var mpPool map[string]*Pool = make(map[string]*Pool, 4)
var lkMp *sync.RWMutex = new(sync.RWMutex)

func GetNoCreate(poolKey string, id interface{}) (IFObject, error) {
	ptPool, ok := GetPool(poolKey)
	if !ok {
		return nil, fmt.Errorf("NOPOOL")
	}
	result, ok := ptPool.GetNoCreate(id)
	if !ok {
		return nil, fmt.Errorf("NULL")
	}
	return result, nil
}

func Get(poolKey string, id interface{}) (IFObject, error) {
	if ptPool, ok := GetPool(poolKey); ok {
		return ptPool.Get(id)
	} else {
		return nil, fmt.Errorf("NOPOOL")
	}
}

func GetPool(poolKey string) (*Pool, bool) {
	lkMp.RLock()
	ptPool, ok := mpPool[poolKey]
	lkMp.RUnlock()
	return ptPool, ok
}

func RegPool(poolKey string, fn CreateFunc) error {
	lkMp.Lock()
	defer lkMp.Unlock()
	if _, ok := mpPool[poolKey]; ok {
		return fmt.Errorf("is exist")
	}
	mpPool[poolKey] = NewPool(fn)
	return nil
}
