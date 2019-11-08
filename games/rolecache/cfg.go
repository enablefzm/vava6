package rolecache

func newCfg() *Cfg {
	pt := &Cfg{
		lowLifeTime:     60,   // 当玩家数低时可以在没有响应该的情况下存在3分钟 180
		normalLifeTime:  60,   //  标准时间是2分钟
		heightLifeTime:  60,   // 如果玩家数量高时只有1分钟
		lowLifeValue:    100,  // 当只有100名在线时时执行
		normalLifeValue: 1000, // 当只有1000名在线时执行
		heartBeat:       60,
	}
	return pt
}

type Cfg struct {
	lowLifeTime     int   // 当数量低时多久释放时间
	normalLifeTime  int   // 标准情况下需要多久释放时间
	heightLifeTime  int   // 当人数高时需要多久释放时间
	lowLifeValue    int   // 当低于这个值时进行慢一点的时间释放
	normalLifeValue int   // 当低于这个值时进行标准时间释放
	heartBeat       int64 // 心跳时间间隔
}
