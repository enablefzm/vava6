package rolecache

import (
	"fmt"
	"sync"
	"time"
	"vava6/vatools"
)

func NewRoleCache() *RoleCache {
	return &RoleCache{
		mpid:   make(map[uint]*Rolec, 1000),
		mpuid:  make(map[string]*Rolec, 1000),
		lk:     new(sync.RWMutex),
		lkLog:  new(sync.Mutex),
		chPool: make(chan bool, 100),
		ptCfg:  newCfg(),
	}
}

// 游戏角色缓存管理
type RoleCache struct {
	mpid    map[uint]*Rolec   // 通过ID的KEY来存放角色对象的MAP
	mpuid   map[string]*Rolec // 通过UID的KEY来存放角色对象的MAP
	lk      *sync.RWMutex     // lk写入数据里的锁
	lkLog   *sync.Mutex       // 记录操作成功失败记数的锁
	chPool  chan bool         // 保存数据的通道
	ptCfg   *Cfg              // 当前缓存的配置信息
	onceRun sync.Once         // 只能运行一次
}

// 开始运行守护精灵
func (this *RoleCache) Start() {
	go func() {
		this.onceRun.Do(this._run)
	}()
}

func (this *RoleCache) _run() {
	ticker := time.NewTicker(time.Duration(this.ptCfg.heartBeat) * time.Second)
	for {
		<-ticker.C
		// valog.OBLog.LogMessage(fmt.Sprintf("开始执行释放空间操作：当前有 %d 个角色", this.GetCount()))
		this.Free()
		// valog.OBLog.LogMessage(fmt.Sprintf("执行完毕：当前有 %d 个角色", this.GetCount()))
	}
}

// 将角色对象放入缓存池子里
//	@parames
//		ob IFRole 角色对象
func (this *RoleCache) Put(ob IFRole) {
	this.lk.Lock()
	this._put(ob)
	this.lk.Unlock()
}

func (this *RoleCache) _put(ob IFRole) {
	p := &Rolec{
		IFRole:   ob,
		lastTime: time.Now().Unix(),
		isSave:   false,
	}
	this.mpid[ob.GetID()] = p
	this.mpuid[ob.GetUID()] = p
}

// 通过ID来获取角色对象
//	@parame
//		id uint 角色ID
//	@return
//		IFRole 	角色对象
//		bool	如果角色不存在则返回false
func (this *RoleCache) GetID(id uint) (IFRole, bool) {
	this.lk.RLock()
	ob, ok := this.mpid[id]
	this.lk.RUnlock()
	if ok {
		return ob.GetRole(), ok
	}
	return nil, ok
}

// 通过UID来获取角色对象
//	@parame
//		uid string 角色UID
//	@return
//		IFRole 	角色对象
//		bool	如果角色不存在则返回false
func (this *RoleCache) GetUID(uid string) (IFRole, bool) {
	this.lk.RLock()
	ob, ok := this.mpuid[uid]
	this.lk.RUnlock()
	if ok {
		return ob.GetRole(), ok
	}
	return nil, ok
}

// 获取当前缓存里有多少个角色对象
func (this *RoleCache) GetCount() int {
	return len(this.mpid)
}

// 将缓存里的数据保存到数据库里
//	@parames
//		lt int64 在缓存里间隔时间，间隔时间大于这个设定时间进行保存
func (this *RoleCache) SaveToDB(lt int64) (count, ok, bad int) {
	// 等待组
	wg := new(sync.WaitGroup)
	// 获取当前的时间
	nt := time.Now().Unix()
	// 获取当前总数量
	count = this.GetCount()
	// 启用读锁标记
	this.lk.RLock()
	// 开始遍历当前所有的缓存对象
	for _, v := range this.mpid {
		// 如果间隔时间大于设定时间且保存标志为未保存时
		if (nt-v.lastTime) > lt && !v.isSave {
			// 保证队列不能超出设定的阀值
			this.chPool <- true
			wg.Add(1)
			go func(ob *Rolec) {
				err := ob.Save()
				// 取出一个占位符
				this.lkLog.Lock()
				if err != nil {
					bad++
					// DEBUG 测试当前错误信息
					fmt.Println("SaveToDB:", err.Error())
				} else {
					ok++
				}
				this.lkLog.Unlock()
				<-this.chPool
				wg.Done()
			}(v)
		}
	}
	this.lk.RUnlock()
	wg.Wait()
	return
}

// 释放缓存空间
//	@return int 当前被释放出来的数量
func (this *RoleCache) Free() int {
	if this.GetCount() < 1 {
		return 0
	}
	// 获取角色无响应存在时间长度
	var lt int64 = int64(this.ptCfg.normalLifeTime)
	var countRole int = this.GetCount()
	switch {
	case countRole < this.ptCfg.lowLifeValue:
		lt = int64(this.ptCfg.lowLifeTime)
	case countRole > this.ptCfg.normalLifeValue:
		lt = int64(this.ptCfg.heightLifeTime)
	}
	// 执行保存
	this.SaveToDB(lt)
	// 获取当前时间
	nt := time.Now().Unix()
	// 定义被释放的数量
	var freeCount int = 0
	// 执行释放
	this.lk.Lock()
	for k, v := range this.mpid {
		if (nt-v.lastTime) > lt && v.isSave {
			delete(this.mpid, k)
			delete(this.mpuid, v.GetUID())
			freeCount++
		}
	}
	this.lk.Unlock()
	return freeCount
}

// 释放全部角色对象，如果没有保存则保存数据
//	@return int 当前被释放的数量
func (this *RoleCache) FreeAll() int {
	this.lk.Lock()
	var freeCount int
	for k, v := range this.mpid {
		if !v.isSave {
			v.Save()
		}
		delete(this.mpid, k)
		delete(this.mpuid, v.GetUID())
		freeCount++
	}
	this.lk.Unlock()
	return freeCount
}

// 获取当前池子里的总共有多少玩家角色
//	@return
//		[]map[string]interface{}
func (this *RoleCache) ListRoles() []map[string]interface{} {
	res := make([]map[string]interface{}, 0, 10)
	maxValue := 0
	this.lk.RLock()
	for _, v := range this.mpid {
		res = append(res, map[string]interface{}{
			"id":       v.GetID(),
			"uid":      v.GetUID(),
			"lastTime": vatools.GetTimeString(v.lastTime),
		})
		maxValue++
		if maxValue >= 10 {
			break
		}
	}
	this.lk.RUnlock()
	return res
}
