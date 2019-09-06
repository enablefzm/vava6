package vatools

import (
	"testing"
)

func TestJson(t *testing.T) {
	mpJson := map[string]interface{}{
		"name": "jimmy",
		"age":  39,
		"sex":  true,
	}

	t.Log(MapToJson(mpJson))

	strJson := `{"age":39,"name":"jimmy","sex":true}`
	mpDb := &{
		
	}

	t.Log(mpDb)
}
