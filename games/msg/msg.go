package msg

import (
	"encoding/json"
)

func fInfo(resMap map[string]interface{}) string {
	v, _ := json.Marshal(resMap)
	return string(v)
}

func NewResMessage(cmd string) *ResMessage {
	obRes := &ResMessage{
		RES: make(map[string]interface{}),
	}
	obRes.addInfo("CMD", cmd)
	return obRes
}

// 返回信息对象
type ResMessage struct {
	RES map[string]interface{}
}

func (t *ResMessage) SetInfo(val interface{}) {
	t.addInfo("INFO", val)
}

func (t *ResMessage) SetRes(blnRes bool, msg string) {
	t.addInfo("RES", blnRes)
	t.addInfo("MSG", msg)
}

func (t *ResMessage) addInfo(key string, val interface{}) {
	t.RES[key] = val
}

func (t *ResMessage) AddInfo(key string, val interface{}) {
	t.addInfo(key, val)
}

func (t *ResMessage) AddInfoValue(key string, val interface{}) {
	if element, ok := t.RES["INFO"]; ok {
		if tInfo, ok := element.(map[string]interface{}); ok {
			tInfo[key] = val
		}
	} else {
		t.RES["INFO"] = map[string]interface{}{
			key: val,
		}
	}
}

func (t *ResMessage) GetRes() map[string]interface{} {
	return t.RES
}

func (t *ResMessage) GetString() string {
	return fInfo(t.RES)
}
