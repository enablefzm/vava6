package servers

import (
	"vava6/vaconn"
)

// 游戏接口对象
type Game interface {
	// 将游戏服务器设定
	SetGameServer(gs *GameServer)
	// 创建玩家信息
	CreatePlayer(conn vaconn.MConn) Player
	// 收到消息处理
	DoCmd(p Player, cmd string)
	// 游戏保存
	Save() error
}

// 玩家接口对象
type Player interface {
	// 获取连接对象
	GetCONN() vaconn.MConn
	// 获取ID
	GetID() uint
	// 获取UID
	GetUID() string
	// 断开事件处理
	HandleConnClose()
}
