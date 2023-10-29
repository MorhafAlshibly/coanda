package pkg

import (
	"github.com/bytedance/sonic"
)

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

func MapStringAnyToMapStringString(input map[string]interface{}) (map[string]string, error) {
	output := make(map[string]string)
	for key, value := range input {
		marshalled, err := sonic.Marshal(value)
		if err != nil {
			return nil, err
		}
		output[key] = string(marshalled)
	}
	return output, nil
}

func MapStringStringToMapStringAny(input map[string]string) map[string]interface{} {
	output := make(map[string]interface{})
	for key, value := range input {
		output[key] = value
	}
	return output
}

func ParsePagination(max *uint32, page *uint64, defaultMax uint8, maxMax uint8) (uint8, uint64) {
	newMax := defaultMax
	newPage := uint64(1)
	if max != nil {
		newMax = uint8(*max)
		if newMax == 0 {
			newMax = defaultMax
		}
		if newMax > maxMax {
			newMax = maxMax
		}
	}
	newPage = 1
	if page != nil {
		if *page > newPage {
			newPage = *page
		}
	}
	return newMax, newPage
}
