package pkg

func RemoveDuplicate[T string | int | uint64](sliceList []T) []T {
	allKeys := make(map[T]bool)
	list := []T{}
	for _, item := range sliceList {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func Contains[T string | int | uint64](sliceList []T, item T) bool {
	for _, sliceItem := range sliceList {
		if sliceItem == item {
			return true
		}
	}
	return false
}

func Remove[T string | int | uint64](sliceList []T, item T) []T {
	list := []T{}
	for _, sliceItem := range sliceList {
		if sliceItem != item {
			list = append(list, sliceItem)
		}
	}
	return list
}
