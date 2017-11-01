package vatools

// 可以定长的数组
type VaArray struct {
	max  uint16
	arrs []interface{}
	idx  uint16
}

func NewVaArray(max uint16) *VaArray {
	return &VaArray{
		max:  max,
		arrs: make([]interface{}, 0, max),
		idx:  0,
	}
}

// 添加数据到定长数组里
func (this *VaArray) Add(v interface{}) {
	// 判断当前日志是否超出长度
	if uint16(len(this.arrs)) >= this.max {
		this.arrs[this.idx] = v
		this.idx++
		if this.idx >= this.max {
			this.idx = 0
		}
	} else {
		this.arrs = append(this.arrs, v)
	}
}

func (this *VaArray) Get() []interface{} {
	if uint16(len(this.arrs)) >= this.max {
		return append(this.arrs[this.idx:], this.arrs[:this.idx]...)
	} else {
		return append(this.arrs)
	}
}
