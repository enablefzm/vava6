package slggame

var ptrManagePlayer *ManagePlayer

func init() {
	ptrManagePlayer = &ManagePlayer{
		loginPlayer:  make(map[uint]*Player, 5000),
		chAddPlayer:  make(chan *Player, 0),
		chMovePlayer: make(chan *Player, 0),
		chFindPlayer: make(chan FindPlayer, 0),
	}
	go ptrManagePlayer.run()
}

func GetFactoryPlayer() *ManagePlayer {
	return ptrManagePlayer
}

type FindPlayer struct {
	id uint
	ch chan *Player
}

type ManagePlayer struct {
	loginPlayer  map[uint]*Player
	chAddPlayer  chan *Player
	chMovePlayer chan *Player
	chFindPlayer chan FindPlayer
}

func (this *ManagePlayer) run() {
	for {
		select {
		case p := <-this.chAddPlayer:
			this.loginPlayer[p.GetID()] = p
		case p := <-this.chMovePlayer:
			if op, ok := this.loginPlayer[p.GetID()]; ok {
				if op.GetCONN() == p.GetCONN() {
					delete(this.loginPlayer, p.GetID())
				}
			}
		}
	}
}

func (this *ManagePlayer) Find(id uint) (*Player, bool) {
	f := FindPlayer{
		id: id,
		ch: make(chan *Player, 0),
	}
	this.chFindPlayer <- f
	p := <-f.ch
	if p != nil {
		return p, true
	} else {
		return nil, false
	}
}
