package vatools

import (
	"testing"
)

func TestProportion(t *testing.T) {
	mp := map[interface{}]int{
		"NULL":      5000,
		"JIMMY1000": 100,
		"Jack100":   200,
		"Egg":       500,
	}

	mp = map[interface{}]int{}

	ptrProportion := NewBaseProportion(1000, mp)
	for i := 0; i < 10; i++ {
		t.Log(ptrProportion.GetRndKey())
		//		rndKey := ptrProportion.GetRndKey()
		//		if rndKey != "JIMMY1000" {
		//			t.Log(rndKey)
		//		}
	}

	t.Log("Number:", SUint8("-12"))

	//	ptrProportion := NewRangeProportion(1000, mp)
	//	iEnd := 0
	//	for {
	//		vKey := ptrProportion.GetRangeKeys(20)
	//		t.Log(vKey)
	//		iEnd++
	//		if iEnd >= 10 {
	//			break
	//		}
	//	}
}
