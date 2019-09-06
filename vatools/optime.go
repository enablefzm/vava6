package vatools

import (
	"time"
)

// type OpTime struct {
// 	lastTime int64
// }

// func

// 计算过去的时间
//	@return
//		int64 记录最后计算的时间
//		int 获得计算次数
func OpTime(lastTime int64, timeSpeed int64) (int64, int) {
	// 获取当前时间
	nowTime := time.Now().Unix()
	// 计算过去时间
	passTime := nowTime - lastTime
	// 如果达不到需要的时间就不需计算
	if passTime < timeSpeed {
		return lastTime, 0
	}
	// 达到返回剩余时间和计算的次数
	opValue := int(passTime / timeSpeed)
	if opValue < 0 {
		opValue = 0
	}
	// 计算剩余时间
	surplus := passTime % timeSpeed
	resultLastTime := nowTime - surplus
	// 返回结果
	return resultLastTime, opValue
}
