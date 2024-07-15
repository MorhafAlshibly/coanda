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
