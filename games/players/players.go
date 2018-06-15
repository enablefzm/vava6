package players

import (
	"fmt"
	"vava6/games/msg"
	"vava6/games/servers"
)

var ptManagePlayer *ManagePlayer = &ManagePlayer{
	mpLoginPlayer:     make(map[string]servers.Player, 1000),
	mpLoginPlayerOnID: make(map[uint]servers.Player, 1000),
	chAdd:             make(chan servers.Player, 1000),
	chMove:            make(chan servers.Player, 1000),
	chGet:             make(chan *proGet, 1000),
	chGetID:           make(chan *proGetOnID, 1000),
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
		chResult: make(chan servers.Player, 2),
	}
	ptManagePlayer.chGet <- result
	p := <-result.chResult
	if p != nil {
		return p, true
	} else {
		return p, false
	}
}

func GetPlayerOnId(id uint) (servers.Player, bool) {
	result := &proGetOnID{
		id:       id,
		chResult: make(chan servers.Player, 2),
	}
	ptManagePlayer.chGetID <- result
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

type proGetOnID struct {
	id       uint
	chResult chan servers.Player
}

type ManagePlayer struct {
	mpLoginPlayer     map[string]servers.Player
	mpLoginPlayerOnID map[uint]servers.Player
	chAdd             chan servers.Player
	chMove            chan servers.Player
	chGet             chan *proGet
	chGetID           chan *proGetOnID
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
			this.mpLoginPlayerOnID[p.GetID()] = p
		case p := <-this.chMove:
			if oldP, ok := this.mpLoginPlayer[p.GetUID()]; ok {
				if oldP == p {
					delete(this.mpLoginPlayer, p.GetUID())
					delete(this.mpLoginPlayerOnID, p.GetID())
				}
			}
		case pro := <-this.chGet:
			uid := pro.uid
			if p, ok := this.mpLoginPlayer[uid]; ok {
				pro.chResult <- p
			} else {
				pro.chResult <- nil
			}
		case proId := <-this.chGetID:
			if p, ok := this.mpLoginPlayerOnID[proId.id]; ok {
				proId.chResult <- p
			} else {
				proId.chResult <- nil
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
