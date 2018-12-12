package pools

import (
	"time"
)

type PoolObject struct {
	ptObject IFObject
	lastTime int64
}

func (this *PoolObject) Save() error {
	return this.ptObject.Save()
}

func (this *PoolObject) Stamp() {
	this.lastTime = time.Now().Unix()
}

func (this *PoolObject) GetLastTime() int64 {
	return this.lastTime
}
