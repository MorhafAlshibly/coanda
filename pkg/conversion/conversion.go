package conversion

import (
	"encoding/json"

	"github.com/MorhafAlshibly/coanda/api"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
)

func Enum[T ~string, PB ~int32](val T, pbmap map[string]int32, dft PB) PB {
	v, ok := pbmap[string(val)]
	if !ok {
		return dft
	}
	return PB(v)
}

func MapToProtobufStruct(m map[string]interface{}) (*structpb.Struct, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	s := &structpb.Struct{}
	err = protojson.Unmarshal(b, s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func ProtobufStructToMap(s *structpb.Struct) (map[string]interface{}, error) {
	b, err := protojson.Marshal(s)
	if err != nil {
		return nil, err
	}
	m := make(map[string]interface{})
	err = json.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func RawJsonToProtobufStruct(m json.RawMessage) (*structpb.Struct, error) {
	s := &structpb.Struct{}
	err := protojson.Unmarshal(m, s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func ProtobufStructToRawJson(s *structpb.Struct) (json.RawMessage, error) {
	b, err := protojson.Marshal(s)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func PaginationToLimitOffset(pagination *api.Pagination, defaultMax uint8, maxMax uint8) (int32, int32) {
	if pagination == nil {
		return int32(defaultMax), 0
	}
	max := uint8(PointerToValue(pagination.Max, uint32(defaultMax)))
	if max > maxMax {
		max = maxMax
	}
	page := PointerToValue(pagination.Page, 1)
	offset := (page - 1) * uint64(max)
	return int32(max), int32(offset)
}

// Pointer to Value
func PointerToValue[T any](p *T, defaultValue T) T {
	if p == nil {
		return defaultValue
	}
	return *p
}

// Value to Pointer
func ValueToPointer[T any](v T) *T {
	return &v
}
