package pools

import (
	"sync"
	"time"
)

const (
	LIVE_TIME int64         = 300
	JUMP_TIME time.Duration = 90
)

type CreateFunc func(id interface{}) (IFObject, error)

type proDel struct {
	key   interface{}
	ptObj *PoolObject
}

func NewPool(fn CreateFunc) *Pool {
	ptPool := &Pool{
		createFun: fn,
		chSave:    make(chan proDel, 10),
		nowTime:   time.Now().Unix(),
	}
	go ptPool.run()
	return ptPool
}

type Pool struct {
	mpPool    sync.Map
	createFun CreateFunc
	chSave    chan proDel
	nowTime   int64
}

// 判断当前数据是否有这个对象
func (this *Pool) Get(id interface{}) (IFObject, error) {
	ptObj, ok := this.GetNoCreate(id)
	if !ok {
		// 创建对象
		var err error
		ptObj, err = this.createFun(id)
		if err != nil {
			return nil, err
		}
		// 创建成功放到到存储对象里
		this.Put(id, ptObj)
	}
	return ptObj, nil
}

func (this *Pool) GetNoCreate(id interface{}) (IFObject, bool) {
	v, ok := this.mpPool.Load(id)
	if ok {
		if res, ok := v.(*PoolObject); ok {
			res.Stamp()
			return res.ptObject, true
		}
	}
	return nil, false
}

func (this *Pool) Put(id interface{}, val IFObject) {
	putVal := &PoolObject{
		ptObject: val,
	}
	putVal.Stamp()
	this.mpPool.Store(id, putVal)
}

// 定时清除没有使用的数据
func (this *Pool) run() {
	go func() {
		for {
			pro := <-this.chSave
			pro.ptObj.Save()
			// 删除当前对象
			if (this.nowTime - pro.ptObj.GetLastTime()) >= LIVE_TIME {
				// 删除
				this.mpPool.Delete(pro.key)
			}
		}
	}()
	// 创建清除时序
	tk := time.NewTicker(JUMP_TIME * time.Second)
	for {
		<-tk.C
		this.nowTime = time.Now().Unix()
		this.mpPool.Range(func(k, v interface{}) bool {
			if r, ok := v.(*PoolObject); ok {
				if (this.nowTime - r.GetLastTime()) >= LIVE_TIME {
					this.chSave <- proDel{
						key:   k,
						ptObj: r,
					}
				}
			} else {
				this.mpPool.Delete(k)
			}
			return true
		})
	}
}

func (this *Pool) Save() error {
	this.mpPool.Range(func(k, v interface{}) bool {
		if r, ok := v.(*PoolObject); ok {
			if err := r.Save(); err != nil {
				// fmt.Println("保存Object出错：", err.Error())
			}
		}
		return true
	})
	return nil
}
