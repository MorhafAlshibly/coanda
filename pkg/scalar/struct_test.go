package scalar

import (
	"testing"

	"github.com/99designs/gqlgen/graphql"
	"google.golang.org/protobuf/types/known/structpb"
)

func Test_MarshalProtobufStruct_StructNil_ReturnGraphQLNull(t *testing.T) {
	m := MarshalProtobufStruct(nil)
	if m != graphql.Null {
		t.Error("Expected nil")
	}
}

func Test_MarshalProtobufStruct_StructNotNil_ReturnGraphQLStruct(t *testing.T) {
	m := MarshalProtobufStruct(&structpb.Struct{})
	if m == nil {
		t.Error("Expected struct")
	}
}

func Test_MarshalProtobufStruct_StructWithFields_ReturnGraphQLStruct(t *testing.T) {
	m := MarshalProtobufStruct(&structpb.Struct{
		Fields: map[string]*structpb.Value{
			"key": {
				Kind: &structpb.Value_StringValue{
					StringValue: "value",
				},
			},
		},
	})
	if m == nil {
		t.Error("Expected struct")
	}
}

func Test_MarshalProtobufStruct_StructWithNestedFields_ReturnGraphQLStruct(t *testing.T) {
	m := MarshalProtobufStruct(&structpb.Struct{
		Fields: map[string]*structpb.Value{
			"key": {
				Kind: &structpb.Value_StructValue{
					StructValue: &structpb.Struct{
						Fields: map[string]*structpb.Value{
							"key": {
								Kind: &structpb.Value_StringValue{
									StringValue: "value",
								},
							},
						},
					},
				},
			},
		},
	})
	if m == nil {
		t.Error("Expected struct")
	}
}

func Test_UnmarshalProtobufStruct_NotMap_ReturnError(t *testing.T) {
	_, err := UnmarshalProtobufStruct("string")
	if err == nil {
		t.Error("Expected error")
	}
}

func Test_UnmarshalProtobufStruct_Map_ReturnStruct(t *testing.T) {
	s, err := UnmarshalProtobufStruct(map[string]interface{}{})
	if err != nil {
		t.Error("Expected nil")
	}
	if s == nil {
		t.Error("Expected struct")
	}
}

func Test_UnmarshalProtobufStruct_MapWithFields_ReturnStruct(t *testing.T) {
	s, err := UnmarshalProtobufStruct(map[string]interface{}{
		"key": "value",
	})
	if err != nil {
		t.Error("Expected nil")
	}
	if s == nil {
		t.Error("Expected struct")
	}
}

func Test_UnmarshalProtobufStruct_MapWithNestedFields_ReturnStruct(t *testing.T) {
	s, err := UnmarshalProtobufStruct(map[string]interface{}{
		"key": map[string]interface{}{
			"key": "value",
		},
	})
	if err != nil {
		t.Error("Expected nil")
	}
	if s == nil {
		t.Error("Expected struct")
	}
}

func Test_UnmarshalProtobufStruct_InvalidType_ReturnError(t *testing.T) {
	_, err := UnmarshalProtobufStruct(1)
	if err == nil {
		t.Error("Expected error")
	}
}
