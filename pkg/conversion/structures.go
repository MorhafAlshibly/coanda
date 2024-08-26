package conversion

import (
	"encoding/json"
	"fmt"
	"strings"

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

// Header is optional, if nil, it will use the keys of the first map as the header
func ArrayOfMapsToCSV(a []map[string]interface{}, header *string) (string, error) {
	if len(a) == 0 {
		return "", nil
	}
	if header == nil {
		keys := make([]string, 0, len(a[0]))
		for k := range a[0] {
			keys = append(keys, k)
		}
		h := strings.Join(keys, ",")
		header = &h
	}
	lines := make([]string, 0, len(a)+1)
	lines = append(lines, *header)
	for _, m := range a {
		line := make([]string, 0, len(m))
		for _, v := range m {
			line = append(line, fmt.Sprintf("%v", v))
		}
		lines = append(lines, strings.Join(line, ","))
	}
	return strings.Join(lines, "\n"), nil
}
