package players

import (
	"fmt"
	"sync"
	"vava6/games/msg"
	"vava6/games/servers"
)

var mpLoginPlayer map[string]servers.Player = make(map[string]servers.Player, 5000)
var lk *sync.RWMutex = new(sync.RWMutex)

func AddPlayer(p servers.Player) {
	lk.Lock()
	if oldP, ok := mpLoginPlayer[p.GetUID()]; ok {
		res := msg.NewResMessage("LOGIN_OTHER")
		res.SetInfo("其它地方登入")
		oldP.GetCONN().Send(res.GetString())
		oldP.GetCONN().Close()
		// 测试关闭状态发送信息是什么样的
		oldP.GetCONN().Send("TEST CLOSE")
	}
	mpLoginPlayer[p.GetUID()] = p
	lk.Unlock()
	showCount()
}

func MovePlayer(p servers.Player) {
	if oldP, ok := mpLoginPlayer[p.GetUID()]; ok {
		lk.Lock()
		if oldP == p {
			delete(mpLoginPlayer, p.GetUID())
		}
		lk.Unlock()
	}
	showCount()
}

func GetPlayerOnUid(uid string) (servers.Player, bool) {
	lk.RLock()
	p, ok := mpLoginPlayer[uid]
	lk.RUnlock()
	return p, ok
}

func Count() int {
	return len(mpLoginPlayer)
}

func showCount() {
	fmt.Println("Players 当前玩家数量", Count())
}
