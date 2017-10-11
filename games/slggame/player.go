package slggame

import (
	"fmt"
	"vava6/vaconn"
)

func NewPlayer(conn vaconn.MConn) *Player {
	return &Player{
		conn: conn,
	}
}

type Player struct {
	id      uint
	conn    vaconn.MConn
	isLogin bool
}

func (this *Player) IsLogin() bool {
	return this.isLogin
}

func (this *Player) SetIsLogin(v bool) {
	this.isLogin = v
}

func (this *Player) GetCONN() vaconn.MConn {
	return this.conn
}

func (this *Player) Send(msg string) error {
	return this.conn.Send(msg)
}

func (this *Player) GetID() uint {
	return this.id
}

func (this *Player) SetID(id uint) {
	this.id = id
}

func (this *Player) HandleConnClose() {
	fmt.Println("slggame player 断开连接处理")
	obGame.DisconnectPlayer(this)
}
