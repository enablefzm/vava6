package vatools

import (
	"testing"
)

func TestProportion(t *testing.T) {
	mp := map[string]int{
		"NULL":      5000,
		"JIMMY1000": 100,
		"Jack100":   200,
		"Egg":       500,
	}

	ptrProportion := NewProportion(1000, mp)
	for i := 0; i < 10; i++ {
		t.Log(ptrProportion.GetRndKey())
		//		rndKey := ptrProportion.GetRndKey()
		//		if rndKey != "JIMMY1000" {
		//			t.Log(rndKey)
		//		}
	}
}
