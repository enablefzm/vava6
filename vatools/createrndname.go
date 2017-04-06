package vatools

var OBCreateName = &createName{
	mpName:    []rune("范林张庄洪周陈丁蒋顾"),
	mpManName: []rune("志旭亮强明武艺富汉泉宇祺恩宸英立鹏勇源"),
	mpWomName: []rune("丽少萍娟玲瑛琪燕秋爱美仙玫玉"),
}

type createName struct {
	mpName    []rune
	mpManName []rune
	mpWomName []rune
}

func (this *createName) GetName() string {
	aLen := len(this.mpName)
	nLen := len(this.mpManName)
	rndVal := CRnd(1, aLen) - 1
	sName1 := make([]rune, 0, 3)
	sName1 = append(sName1, this.mpName[rndVal])
	rndHow := CRnd(1, 100)
	il := 3
	switch {
	case rndHow < 30:
		il = 1
	default:
		il = 2
	}
	for i := 0; i < il; i++ {
		vRnd := CRnd(1, nLen) - 1
		sName1 = append(sName1, this.mpManName[vRnd])
	}
	return string(sName1)
}
