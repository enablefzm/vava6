package vatools

import (
	"errors"
	"fmt"
)

type BaseProportion struct {
	maxProportion int
	mpValue       map[interface{}]proportionValue
}

func NewBaseProportion(iMax int, mpProValue map[interface{}]int) *BaseProportion {
	pt := &BaseProportion{
		maxProportion: iMax,
	}
	pt.InitValue(mpProValue)
	return pt
}

func (this *BaseProportion) InitValue(mpProValue map[interface{}]int) {
	var iCount int
	for _, v := range mpProValue {
		iCount += v
	}
	var dCount float64 = float64(iCount)
	var dMax float64 = float64(this.maxProportion)
	var tMin uint16 = 0
	// 计算占比
	this.mpValue = make(map[interface{}]proportionValue, len(mpProValue))
	for k, v := range mpProValue {
		r := float64(v) / dCount * dMax
		tMax := uint16(r) + tMin
		this.mpValue[k] = proportionValue{
			min: tMin,
			max: tMax,
		}
		tMin = tMax
	}
}

func (this *BaseProportion) GetRndKey() interface{} {
	intVal := uint16(CRnd(0, this.maxProportion))
	var result interface{}
	for k, v := range this.mpValue {
		result = k
		if intVal >= v.min && intVal <= v.max {
			return k
		}
	}
	return result
}

// 比例分布运算
type Proportion struct {
	maxProportion uint16                     // 占比最大值
	mpValue       map[string]proportionValue // 占比明细
}

// 构造比例分布占比
func NewProportion(iMax uint16, mpProValue map[string]int) *Proportion {
	pt := &Proportion{
		maxProportion: iMax,
	}
	pt.InitValue(mpProValue)
	return pt
}

func (this *Proportion) InitValue(mpProValue map[string]int) {
	// 计算总值
	var iCount int
	for _, v := range mpProValue {
		iCount += v
	}
	var dCount float64 = float64(iCount)
	var dMax float64 = float64(this.maxProportion)
	var tMin uint16 = 0
	// 计算占比
	this.mpValue = make(map[string]proportionValue, len(mpProValue))
	for k, v := range mpProValue {
		r := float64(v) / dCount * dMax
		tMax := uint16(r) + tMin
		this.mpValue[k] = proportionValue{
			min: tMin,
			max: tMax,
		}
		tMin = tMax
	}
}

func (this *Proportion) GetInfo() map[string]string {
	res := make(map[string]string, len(this.mpValue))
	for k, v := range this.mpValue {
		res[k] = fmt.Sprint(k, " : ", v.min, " - ", v.max)
	}
	return res
}

// 通过占比值返回随机Key
func (this *Proportion) GetRndKey() string {
	intVal := uint16(CRnd(0, int(this.maxProportion)))
	var result string
	for k, v := range this.mpValue {
		result = k
		if intVal >= v.min && intVal <= v.max {
			return k
		}
	}
	return result
}

type proportionValue struct {
	min uint16
	max uint16
}

// 创建可以被获取的分布运算
type CanGetProportion struct {
	pro   *Proportion
	value map[string]int
}

func NewCanGetProportion(iMax uint16, mpProValue map[string]int) *CanGetProportion {
	pt := NewProportion(iMax, mpProValue)
	return &CanGetProportion{
		pro:   pt,
		value: mpProValue,
	}
}

// 一直获取数据
func (this *CanGetProportion) GetValue() (string, error) {
	k := this.pro.GetRndKey()
	if this.value[k] < 1 {
		delete(this.value, k)
		// 判断是否还有数量
		for k, v := range this.value {
			if v < 1 {
				delete(this.value, k)
			}
		}
		if len(this.value) < 1 {
			return "", errors.New("NULL")
		}
		// 重新分配对象
		this.pro.InitValue(this.value)
		return this.GetValue()
	}
	this.value[k]--
	return k, nil
}
