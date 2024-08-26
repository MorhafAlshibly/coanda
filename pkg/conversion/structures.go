package conversion

import (
	"encoding/json"
	"fmt"
	"sort"
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

func ArrayOfMapsToCsv(a []map[string]interface{}) (string, error) {
	if len(a) == 0 {
		return "", nil
	}
	keys := make([]string, 0, len(a[0]))
	for k := range a[0] {
		keys = append(keys, k)
	}
	// Sort the keys to ensure consistent order
	sort.Strings(keys)
	h := strings.Join(keys, ",")
	lines := make([]string, 0, len(a)+1)
	lines = append(lines, h)
	for _, m := range a {
		// Prepare the line so we can add the values in the correct order
		line := make([]string, len(m))
		for k, v := range m {
			value := fmt.Sprintf("%v", v)
			// If the value is a json raw message, convert it to a string
			if _, ok := v.(json.RawMessage); ok {
				value = string(v.(json.RawMessage))
			}
			// If the value contains a comma, wrap it in double quotes and escape any double quotes
			if strings.Contains(value, ",") {
				value = fmt.Sprintf("\"%v\"", strings.ReplaceAll(value, "\"", "\"\""))
			}
			// Place the value in the correct column
			lineIndex := sort.SearchStrings(keys, k)
			line[lineIndex] = value
		}
		lines = append(lines, strings.Join(line, ","))
	}
	return strings.Join(lines, "\n"), nil
}
