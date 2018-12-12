package pools

import (
	"fmt"
	"testing"
	"time"
)

// func(id interface{}) (IFObject, error)
func TestPool(t *testing.T) {
	RegPool("TST", func(id interface{}) (IFObject, error) {
		fmt.Println("创建了", id)
		return &tstObject{
			id: id,
		}, nil
	})

	v, err := GetNoCreate("TST", 1)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("操作完成", v)
	}

	fmt.Println(Get("TST", 2))
	fmt.Println(Get("TST", 1))

	go func() {
		fmt.Println("")
		for {
			time.Sleep(1 * time.Second)
			fmt.Print(".")
		}
	}()

	fmt.Println("sleep....")
	time.Sleep(19 * time.Second)

	fmt.Println(Get("TST", 1))
	fmt.Println(Get("TST", 2))
	fmt.Println(Get("TST", 3))
}

type tstObject struct {
	id interface{}
}

func (this *tstObject) Save() error {
	fmt.Println(this.id, " Saveing...")
	return nil
}
