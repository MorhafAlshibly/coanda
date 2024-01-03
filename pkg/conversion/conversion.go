package conversion

import (
	"encoding/json"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
)

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

func PageToOffset(page *uint64, max uint8) int32 {
	if page == nil {
		return 0
	}
	return int32((*page - 1) * uint64(max))
}
