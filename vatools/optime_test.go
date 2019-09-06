package vatools

import (
	"testing"
	"time"
)

func TestOpTime(t *testing.T) {
	var speedTime int64 = 10
	nowTime := time.Now().Unix()
	lastTime := nowTime - 10
	t.Log("now Time:", nowTime)
	t.Log(OpTime(lastTime, speedTime))
}
