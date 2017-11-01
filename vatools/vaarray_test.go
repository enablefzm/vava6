package vatools

import (
	"fmt"
	"testing"
	"time"
)

func TestVaArray(t *testing.T) {
	arr := NewVaArray(10)
	for i := 0; i < 16; i++ {
		arr.Add(i)
	}
	v1 := arr.Get()
	v2 := arr.Get()
	t.Log(v1, v2)
	fmt.Printf("Point %p %p \n", v1, v2)

	if v, ok := v1[0].(int); ok {
		fmt.Println(v, "is []int")
	} else {
		fmt.Println(v, "is not []int")
	}
	tt := time.Now().Format("1/2 15:04:05")
	fmt.Println(tt)
}
