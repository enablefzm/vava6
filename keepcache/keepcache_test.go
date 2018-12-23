package keepcache

import (
	"fmt"
	"sync"
	"testing"
)

func TestKeepCache(t *testing.T) {
	kc := NewKeepCache(func(k interface{}) (interface{}, error) {
		fmt.Println("创建了", k)
		return k, nil
	})

	wg := new(sync.WaitGroup)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			kc.Get(1)
			kc.Get(2)
			fmt.Println("OK")
			wg.Done()
		}()
	}
	wg.Wait()
}
