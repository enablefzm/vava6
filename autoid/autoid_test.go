package autoid

import (
	"fmt"
	"sync"
	"testing"
	"vava6/vatools"
)

func TestAutoID(t *testing.T) {
	db := &AutoIdDB{}
	err := vatools.UnJson(`{"LastID":5,"CacheID":[4,3,2,1,5]}`, db)
	if err != nil {
		fmt.Println("UnJson Err: ", err.Error())
	}
	pAutoID := NewAutoIDOnDb(db)
	fmt.Println(pAutoID.GetJsonSave())
	wt := new(sync.WaitGroup)
	for i := 0; i < 20; i++ {
		// fmt.Println(pAutoID.GetID())
		wt.Add(1)
		go func() {
			id := pAutoID.GetID()
			fmt.Println(id)
			pAutoID.PutID(id)
			wt.Done()
		}()
	}
	//	pAutoID.PutID(2)
	//	fmt.Println(pAutoID.GetID())
	wt.Wait()
	fmt.Println(pAutoID.GetJsonSave())
}
