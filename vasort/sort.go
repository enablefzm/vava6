package vasort

// 冒泡的排序接口
type IFSortBubbling interface {
	Len() int
	Swap(i, j int)
	Less(i, j int) bool
}

// 执行冒泡排序
func SortPubbling(obArr IFSortBubbling) {
	il := obArr.Len()
	for i := il - 1; i > 0; i-- {
		j := i - 1
		if j < 0 {
			break
		}
		if obArr.Less(i, j) {
			obArr.Swap(i, j)
		} else {
			break
		}
	}
}
