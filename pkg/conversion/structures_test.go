package conversion

import (
	"encoding/json"
	"testing"

	"google.golang.org/protobuf/types/known/structpb"
)

func Test_MapToProtobufStruct_EmptyMap_EmptyStruct(t *testing.T) {
	m := map[string]interface{}{}
	actual, err := MapToProtobufStruct(m)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(actual.Fields) != 0 {
		t.Errorf("Expected %v but got %v", 0, len(actual.Fields))
	}
}

func Test_MapToProtobufStruct_NonEmptyMap_NonEmptyStruct(t *testing.T) {
	m := map[string]interface{}{"a": 1}
	actual, err := MapToProtobufStruct(m)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(actual.Fields) != 1 {
		t.Errorf("Expected %v but got %v", 1, len(actual.Fields))
	}
}

func Test_ProtobufStructToMap_EmptyStruct_EmptyMap(t *testing.T) {
	s := &structpb.Struct{}
	actual, err := ProtobufStructToMap(s)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(actual) != 0 {
		t.Errorf("Expected %v but got %v", 0, len(actual))
	}
}

func Test_ProtobufStructToMap_NonEmptyStruct_NonEmptyMap(t *testing.T) {
	s := &structpb.Struct{Fields: map[string]*structpb.Value{"a": {Kind: &structpb.Value_NumberValue{NumberValue: 1}}}}
	actual, err := ProtobufStructToMap(s)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(actual) != 1 {
		t.Errorf("Expected %v but got %v", 1, len(actual))
	}
}

func Test_RawJsonToProtobufStruct_EmptyJson_EmptyStruct(t *testing.T) {
	m := json.RawMessage("{}")
	actual, err := RawJsonToProtobufStruct(m)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(actual.Fields) != 0 {
		t.Errorf("Expected %v but got %v", 0, len(actual.Fields))
	}
}

func Test_RawJsonToProtobufStruct_NonEmptyJson_NonEmptyStruct(t *testing.T) {
	m := json.RawMessage(`{"a": 1}`)
	actual, err := RawJsonToProtobufStruct(m)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(actual.Fields) != 1 {
		t.Errorf("Expected %v but got %v", 1, len(actual.Fields))
	}
}

func Test_ProtobufStructToRawJson_EmptyStruct_EmptyJson(t *testing.T) {
	s := &structpb.Struct{}
	actual, err := ProtobufStructToRawJson(s)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(actual) != 2 {
		t.Errorf("Expected %v but got %v", 2, len(actual))
	}
}

func Test_ProtobufStructToRawJson_NonEmptyStruct_NonEmptyJson(t *testing.T) {
	s := &structpb.Struct{Fields: map[string]*structpb.Value{"a": {Kind: &structpb.Value_NumberValue{NumberValue: 1}}}}
	actual, err := ProtobufStructToRawJson(s)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(actual) == 0 {
		t.Errorf("Expected non-empty but got empty")
	}
}

func Test_RawJsonToMap_EmptyJson_EmptyMap(t *testing.T) {
	m := json.RawMessage("{}")
	actual, err := RawJsonToMap(m)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(actual) != 0 {
		t.Errorf("Expected %v but got %v", 0, len(actual))
	}
}

func Test_RawJsonToMap_NonEmptyJson_NonEmptyMap(t *testing.T) {
	m := json.RawMessage(`{"a": 1}`)
	actual, err := RawJsonToMap(m)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(actual) != 1 {
		t.Errorf("Expected %v but got %v", 1, len(actual))
	}
}

func Test_MapToRawJson_EmptyMap_EmptyJson(t *testing.T) {
	m := map[string]interface{}{}
	actual, err := MapToRawJson(m)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(actual) != 2 {
		t.Errorf("Expected %v but got %v", 2, len(actual))
	}
}

func Test_MapToRawJson_NonEmptyMap_NonEmptyJson(t *testing.T) {
	m := map[string]interface{}{"a": 1}
	actual, err := MapToRawJson(m)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(actual) == 0 {
		t.Errorf("Expected non-empty but got empty")
	}
}

func Test_RawJsonToArrayOfMaps_EmptyJson_EmptyArray(t *testing.T) {
	m := json.RawMessage("[]")
	actual, err := RawJsonToArrayOfMaps(m)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(actual) != 0 {
		t.Errorf("Expected %v but got %v", 0, len(actual))
	}
}

func Test_RawJsonToArrayOfMaps_NonEmptyJson_NonEmptyArray(t *testing.T) {
	m := json.RawMessage(`[{"a": 1}]`)
	actual, err := RawJsonToArrayOfMaps(m)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(actual) != 1 {
		t.Errorf("Expected %v but got %v", 1, len(actual))
	}
}

func Test_ArrayOfMapsToCsv_EmptyArray_EmptyCSV(t *testing.T) {
	m := []map[string]interface{}{}
	actual, err := ArrayOfMapsToCsv(m)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(actual) != 0 {
		t.Errorf("Expected %v but got %v", "", len(actual))
	}
}

func Test_ArrayOfMapsToCsv_NonEmptyArray_NonEmptyCSV(t *testing.T) {
	m := []map[string]interface{}{{"ab": 1, "b": 2, "c": "test"}, {"b": 5, "ab": "test1", "c": "test2"}}
	actual, err := ArrayOfMapsToCsv(m)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(actual) == 0 {
		t.Errorf("Expected non-empty but got empty")
	}
	expected := "ab,b,c\n1,2,test\ntest1,5,test2"
	if actual != expected {
		t.Errorf("Expected %v but got %v", expected, actual)
	}
}

func Test_ArrayOfMapsToCsv_WithJsonRawMessage_JsonStringInCsv(t *testing.T) {
	m := []map[string]interface{}{{"ab": 1, "b": 2, "c": json.RawMessage(`{"a": 1, "d": 2}`)}}
	actual, err := ArrayOfMapsToCsv(m)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(actual) == 0 {
		t.Errorf("Expected non-empty but got empty")
	}
	expected := "ab,b,c\n1,2,\"{\"\"a\"\": 1, \"\"d\"\": 2}\""
	if actual != expected {
		t.Errorf("Expected %v but got %v", expected, actual)
	}
}
