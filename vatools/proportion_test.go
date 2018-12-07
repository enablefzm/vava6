package vatools

import (
	"testing"
)

func TestProportion(t *testing.T) {
	arr := GetRndInts(1, 10, 8)
	t.Log(arr)
	for i := 0; i < len(arr); i++ {
		for j := i + 1; j < len(arr); j++ {
			if arr[i] > arr[j] {
				arr[i], arr[j] = arr[j], arr[i]
			}
		}
	}
	t.Log(arr)
}
