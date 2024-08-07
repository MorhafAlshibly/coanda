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

func RawJsonToMap(m json.RawMessage) (map[string]interface{}, error) {
	var i interface{}
	err := json.Unmarshal(m, &i)
	if err != nil {
		return nil, err
	}
	return i.(map[string]interface{}), nil
}

func MapToRawJson(m map[string]interface{}) (json.RawMessage, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func RawJsonToArrayOfMaps(m json.RawMessage) ([]map[string]interface{}, error) {
	var i []interface{}
	err := json.Unmarshal(m, &i)
	if err != nil {
		return nil, err
	}
	a := make([]map[string]interface{}, len(i))
	for j, v := range i {
		a[j] = v.(map[string]interface{})
	}
	return a, nil
}
