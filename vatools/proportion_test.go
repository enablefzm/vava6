package vatools

import (
	"testing"
)

func TestProportion(t *testing.T) {
	for i := 0; i < 10; i++ {
		arr := GetRndInts(1, 15000, 100)
		tSort(arr)
		t.Log(arr)
	}
}

func tSort(arr []int) {
	for i := 0; i < len(arr); i++ {
		for j := i + 1; j < len(arr); j++ {
			if arr[i] > arr[j] {
				arr[i], arr[j] = arr[j], arr[i]
			}
		}
	}
}
