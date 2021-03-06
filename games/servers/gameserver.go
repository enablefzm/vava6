package servers

import (
	"fmt"
	"io"
	"strings"
	"vava6/vaconn"
)

const (
	SER_STOP    = 0
	SER_RUNING  = 1
	SER_PAUSE   = 2
	SER_STOPING = 4
)

func NewGameServer(g Game) *GameServer {
	gs := &GameServer{
		state:        SER_STOP,
		scks:         make([]*vaconn.NaviConnect, 0, 2),
		CHClose:      make(chan bool, 0),
		chAddPlayer:  make(chan Player, 0),
		chMovePlayer: make(chan Player, 0),
		game:         g,
		mpPlayer:     make(map[vaconn.MConn]Player, 5000),
	}
	// 设定当前的游戏服务器和当前加载的游戏进行关联设定
	g.SetGameServer(gs)
	return gs
}

// 游戏服务器
type GameServer struct {
	state        byte                    // 当前游戏有务器状态
	scks         []*vaconn.NaviConnect   // 开放监听的连接对象
	CHClose      chan bool               // 关闭通道
	chAddPlayer  chan Player             // 添加已登入游戏的玩家通道
	chMovePlayer chan Player             // 删除玩家连接对象
	game         Game                    // 游戏对象 - 可以加载不同类型的游戏
	mpPlayer     map[vaconn.MConn]Player // 存放已建连接但是未验证登入的玩家连接对象
	isPause      bool                    // 暂停标记
	isRunStop    bool                    // 正在运行游戏服务器停止
}

// 一次加载所有连接对象
func (this *GameServer) AddSockets(conns ...*vaconn.NaviConnect) {
	for _, conn := range conns {
		this.AddSocket(conn)
	}
}

// 添加新套接字
func (this *GameServer) AddSocket(conn *vaconn.NaviConnect) {
	conn.SetLinkFunc(this.handleNewPlayer)
	// 设定连接对象发生错误时
	conn.SetError(this.onConnError)
	// 添加到当前游戏服务的连接器中
	this.scks = append(this.scks, conn)
}

// 暂停玩家连接
func (this *GameServer) Pause() {
	this.state = SER_PAUSE
	this.isPause = true
}

func (this *GameServer) Resume() {
	this.state = SER_RUNING
	this.isPause = false
}

// 执行停止游戏服务
func (this *GameServer) Stop() {
	// 执行游戏逻辑里的保存操作
	this.CHClose <- true
}

// 延尺执行停止游戏服务器
func (this *GameServer) LastStop(sec int) error {
	if this.isRunStop == true {
		return fmt.Errorf("服务器正在停止中...")
	}
	// 标记上停止
	this.isRunStop = true
	return nil
}

func (this *GameServer) Start() {
	if this.state == SER_RUNING {
		return
	}
	// 将服务器状态标记为运行状态
	this.state = SER_RUNING
	for _, sck := range this.scks {
		go sck.Listen()
	}
	for {
		select {
		// 监听是否有关闭游戏的信息
		case blnClose := <-this.CHClose:
			if blnClose == true {
				close(this.CHClose)
				// 正在停止中
				this.state = SER_STOPING
				// 断开所有的连接
				for _, p := range this.mpPlayer {
					p.GetCONN().Close()
				}
				// 保存游戏逻辑服务器里的保存操作
				err := this.game.Save()
				if err != nil {
					fmt.Println("【GameServer】保存游戏数据出错：", err.Error())
				}
				fmt.Println("【GameServer】Shutdown！")
				// TODO...
				return
			}
		// 放入建立接接但未登入游戏的玩家
		case p := <-this.chAddPlayer:
			this.mpPlayer[p.GetCONN()] = p
			// this.ShowCountPlayer()
		// 移除未登入游戏
		case p := <-this.chMovePlayer:
			delete(this.mpPlayer, p.GetCONN())
			// this.ShowCountPlayer()
		}
	}
}

// 连接对象发生错误时执行
func (this *GameServer) onConnError(err error) {
	fmt.Println(err.Error())
	this.CHClose <- true
}

// 新玩家连接进入游戏
func (this *GameServer) handleNewPlayer(conn vaconn.MConn) {
	defer conn.Close()
	if this.state != SER_RUNING {
		return
	}
	player := this.game.CreatePlayer(conn)
	// fmt.Println(player.GetID(), "连接进入")
	// 添加到未登入玩家对象
	this.AddPlayer(player)
	for {
		strCmd, err := conn.Read()
		if err != nil {
			if err == io.EOF {
				// 正常连接断开
			} else {
				// 非正常断开连接
			}
			// fmt.Println(player.GetID(), "玩家连接被断开了", player.GetCONN().GetIPInfo())
			player.HandleConnClose()
			// 删除当前的连接对象
			this.MovePlayer(player)
			return
		}

		arrCmds := strings.Split(strCmd, "\r\n")
		// 转交给命令对象处理
		for _, cmd := range arrCmds {
			if len(cmd) < 1 {
				continue
			}
			// 转交给游戏处理消息
			this.game.DoCmd(player, cmd)
		}
	}
}

func (this *GameServer) AddPlayer(p Player) {
	this.chAddPlayer <- p
}

func (this *GameServer) MovePlayer(p Player) {
	this.chMovePlayer <- p
}

func (this *GameServer) GetLinkCount() int {
	return len(this.mpPlayer)
}

func (this *GameServer) ShowCountPlayer() {
	fmt.Println("GameServer当前连接进入的玩家：", len(this.mpPlayer))
}
