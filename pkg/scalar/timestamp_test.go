package scalar

import (
	"testing"

	"github.com/99designs/gqlgen/graphql"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Test_MarshalProtobufTimestamp_TimestampNil_ReturnGraphQLNull(t *testing.T) {
	m := MarshalProtobufTimestamp(nil)
	if m != graphql.Null {
		t.Error("Expected nil")
	}
}

func Test_MarshalProtobufTimestamp_TimestampNotNil_ReturnGraphQLTimestamp(t *testing.T) {
	m := MarshalProtobufTimestamp(&timestamppb.Timestamp{})
	if m == nil {
		t.Error("Expected timestamp")
	}
}

func Test_UnmarshalProtobufTimestamp_InvalidTime_ReturnError(t *testing.T) {
	_, err := UnmarshalProtobufTimestamp("invalid")
	if err == nil {
		t.Error("Expected error")
	}
}

func Test_UnmarshalProtobufTimestamp_ValidTime_ReturnTimestamp(t *testing.T) {
	timestamp, err := UnmarshalProtobufTimestamp("2021-01-01T00:00:00Z")
	if err != nil {
		t.Error("Expected nil")
	}
	if timestamp == nil {
		t.Error("Expected timestamp")
	}
}

func Test_UnmarshalProtobufTimestamp_InvalidType_ReturnError(t *testing.T) {
	_, err := UnmarshalProtobufTimestamp(1)
	if err == nil {
		t.Error("Expected error")
	}
}
