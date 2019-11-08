package rolecache

import (
	"time"
)

type IFRole interface {
	GetID() uint
	GetUID() string
	Save() error
}

func NewRolec(obRole IFRole) *Rolec {
	return &Rolec{
		IFRole:   obRole,
		lastTime: time.Now().Unix(),
		isSave:   true,
	}
}

// 被存入进来的角色对象
type Rolec struct {
	IFRole
	lastTime int64
	isSave   bool
}

// 设定被该问到的标志
func (this *Rolec) SetFlag() {
	this.lastTime = time.Now().Unix()
	this.isSave = false
}

// 保存
func (this *Rolec) Save() error {
	err := this.IFRole.Save()
	if err == nil {
		this.isSave = true
	}
	return err
}

// 获取角色对象并且打上被访问的时间标记
//	@return IFRole
func (this *Rolec) GetRole() IFRole {
	this.SetFlag()
	return this.IFRole
}
