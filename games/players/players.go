package players

import (
	"fmt"
	"vava6/games/msg"
	"vava6/games/servers"
)

var ptManagePlayer *ManagePlayer = &ManagePlayer{
	mpLoginPlayer: make(map[string]servers.Player, 1000),
	chAdd:         make(chan servers.Player, 1000),
	chMove:        make(chan servers.Player, 1000),
	chGet:         make(chan *proGet, 1000),
}

func init() {
	go ptManagePlayer.Run()
}

func AddPlayer(p servers.Player) {
	ptManagePlayer.chAdd <- p
}

func MovePlayer(p servers.Player) {
	ptManagePlayer.chMove <- p
}

func GetPlayerOnUid(uid string) (servers.Player, bool) {
	result := &proGet{
		uid:      uid,
		chResult: make(chan servers.Player, 0),
	}
	ptManagePlayer.chGet <- result
	p := <-result.chResult
	if p != nil {
		return p, true
	} else {
		return p, false
	}
}

func Count() int {
	return 0
}

type proGet struct {
	uid      string
	chResult chan servers.Player
}

type ManagePlayer struct {
	mpLoginPlayer map[string]servers.Player
	chAdd         chan servers.Player
	chMove        chan servers.Player
	chGet         chan *proGet
}

func (this *ManagePlayer) Run() {
	for {
		select {
		case p := <-this.chAdd:
			if oldP, ok := this.mpLoginPlayer[p.GetUID()]; ok {
				res := msg.NewResMessage("LOGIN_OTHER")
				res.SetInfo("其它地方登入")
				oldP.GetCONN().Send(res.GetString())
				oldP.GetCONN().Close()
			}
			this.mpLoginPlayer[p.GetUID()] = p
		case p := <-this.chMove:
			if oldP, ok := this.mpLoginPlayer[p.GetUID()]; ok {
				if oldP == p {
					delete(this.mpLoginPlayer, p.GetUID())
				}
			}
		case pro := <-this.chGet:
			uid := pro.uid
			if p, ok := this.mpLoginPlayer[uid]; ok {
				pro.chResult <- p
			} else {
				pro.chResult <- nil
			}
		}
	}
}

func (this *ManagePlayer) showCount() {
	fmt.Println("ManagePlayer当前玩家数量", len(this.mpLoginPlayer))
	for k, v := range this.mpLoginPlayer {
		fmt.Println(k, v.GetCONN().GetIPInfo())
	}
}
