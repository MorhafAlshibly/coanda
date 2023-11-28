package pkg

import (
	"context"

	"github.com/MorhafAlshibly/coanda/api"
	"github.com/bytedance/sonic"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

func MapStringAnyToMapStringString(input map[string]any) (map[string]string, error) {
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

func MapStringStringToMapStringAny(input map[string]string) (map[string]any, error) {
	output := make(map[string]any)
	for key, value := range input {
		var unmarshalled any
		err := sonic.Unmarshal([]byte(value), &unmarshalled)
		if err != nil {
			return nil, err
		}
		output[key] = unmarshalled
	}
	return output, nil
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

func CursorToDocuments[T *api.Record | *api.Item | *api.Team](ctx context.Context, cursor *mongo.Cursor, parseFunc func(*mongo.Cursor) (T, error), page uint64, max uint8) ([]T, error) {
	var result []T
	skip := (int(page) - 1) * int(max)
	for i := 0; i < skip; i++ {
		cursor.Next(ctx)
	}
	for i := 0; i < int(max); i++ {
		if !cursor.Next(ctx) {
			break
		}
		cursorResult, err := parseFunc(cursor)
		if err != nil {
			return nil, err
		}
		result = append(result, cursorResult)
	}
	return result, nil
}

func InterfaceToInt64(input interface{}) int64 {
	switch input.(type) {
	case int:
		return int64(input.(int))
	case int8:
		return int64(input.(int8))
	case int16:
		return int64(input.(int16))
	case int32:
		return int64(input.(int32))
	case int64:
		return input.(int64)
	case uint:
		return int64(input.(uint))
	case uint8:
		return int64(input.(uint8))
	case uint16:
		return int64(input.(uint16))
	case uint32:
		return int64(input.(uint32))
	case uint64:
		return int64(input.(uint64))
	default:
		return 0
	}
}

func InterfaceToUint64(input interface{}) uint64 {
	switch input.(type) {
	case int:
		return uint64(input.(int))
	case int8:
		return uint64(input.(int8))
	case int16:
		return uint64(input.(int16))
	case int32:
		return uint64(input.(int32))
	case int64:
		return uint64(input.(int64))
	case uint:
		return uint64(input.(uint))
	case uint8:
		return uint64(input.(uint8))
	case uint16:
		return uint64(input.(uint16))
	case uint32:
		return uint64(input.(uint32))
	case uint64:
		return input.(uint64)
	default:
		return 0
	}
}

func MapStringAnyToBsonD(input map[string]any) bson.D {
	output := bson.D{}
	for key, value := range input {
		output = append(output, bson.E{Key: key, Value: value})
	}
	return output
}
