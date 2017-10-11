package slggame

import (
	"fmt"
	"strings"
	"vava6/games/servers"
	"vava6/vaconn"
)

// 选通过InitGame构造包的全局游戏对象
// 通过GetGame来获取游戏对象
// 游戏对象
var obGame *Game

// 游戏逻辑对象
type IFGameLogic interface {
	OperateCmd(g *Game, p *Player, cmd string, args []string)
	DisconnectPlayer(p *Player)
	Save() error
}

func InitGame(ptrGame IFGameLogic) *Game {
	obGame = &Game{
		gameLogic: ptrGame,
	}
	return obGame
}

func GetGame() *Game {
	return obGame
}

type Game struct {
	gs        *servers.GameServer
	gameLogic IFGameLogic
}

func (this *Game) SetGameServer(gs *servers.GameServer) {
	this.gs = gs
}

func (this *Game) GetGameServer() *servers.GameServer {
	return this.gs
}

func (this *Game) CreatePlayer(conn vaconn.MConn) servers.Player {
	return NewPlayer(conn)
}

func (this *Game) DoCmd(p servers.Player, cmd string) {
	fmt.Println(p.GetCONN(), " 收到命令 ", cmd)
	if np, ok := p.(*Player); ok {
		arrc := strings.Split(cmd, " ")
		c := arrc[0]
		// 让游戏逻辑对象去处理
		this.gameLogic.OperateCmd(this, np, c, arrc[1:])
	}
}

func (this *Game) Save() error {
	return this.gameLogic.Save()
}

func (this *Game) DisconnectPlayer(p *Player) {
	this.gameLogic.DisconnectPlayer(p)
}
