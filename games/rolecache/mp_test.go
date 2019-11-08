package rolecache

import (
	"fmt"
	"sync"
	"testing"
)

func TestMp(t *testing.T) {
	var mp map[int]int = map[int]int{
		1: 1,
		2: 2,
		3: 3,
		4: 4,
		5: 5,
	}
	for k, v := range mp {
		if k == 2 || k == 3 || k == 5 {
			delete(mp, k)
			fmt.Println("delete ", v)
		}
	}
	t.Log(mp)
	var n1 sync.Once
	n1.Do(func() {
		fmt.Println("n1 do")
	})

	n1.Do(func() {
		fmt.Println("n2 do")
	})
	startLog()
	startLog()
	go func() { startLog() }()
	go func() { startLog() }()
	go func() { startLog() }()

}

func startLog() {
	var ptOnce sync.Once
	ptOnce.Do(logVal)
}

func logVal() {
	fmt.Println("start go")
}
