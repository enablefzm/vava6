package slggame

import (
	"sync"
)

var mpCmds = make(map[string]func(*Game, *Player, string, []string), 100)
var lk = new(sync.Mutex)

func RegCmd(cmdName string, fc func(*Game, *Player, string, []string)) {
	lk.Lock()
	mpCmds[cmdName] = fc
	lk.Unlock()
}

func GetCMD(cmdName string) (func(*Game, *Player, string, []string), bool) {
	fc, ok := mpCmds[cmdName]
	return fc, ok
}
