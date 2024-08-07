package conversion

import (
	"database/sql"
	"testing"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func Test_StringToSqlNullString_NilPointer_EmptyAndInvalid(t *testing.T) {
	actual := StringToSqlNullString(nil)
	if actual.String != "" {
		t.Errorf("Expected %v but got %v", "", actual.String)
	}
	if actual.Valid != false {
		t.Errorf("Expected %v but got %v", false, actual.Valid)
	}
}

func Test_StringToSqlNullString_NotNilPointer_StringAndValid(t *testing.T) {
	s := "test"
	actual := StringToSqlNullString(&s)
	if actual.String != s {
		t.Errorf("Expected %v but got %v", s, actual.String)
	}
	if actual.Valid != true {
		t.Errorf("Expected %v but got %v", true, actual.Valid)
	}
}

func Test_Int64ToSqlNullInt64_NilPointer_ZeroAndInvalid(t *testing.T) {
	actual := Int64ToSqlNullInt64(nil)
	if actual.Int64 != 0 {
		t.Errorf("Expected %v but got %v", 0, actual.Int64)
	}
	if actual.Valid != false {
		t.Errorf("Expected %v but got %v", false, actual.Valid)
	}
}

func Test_Int64ToSqlNullInt64_NotNilPointer_Int64AndValid(t *testing.T) {
	i := int64(10)
	actual := Int64ToSqlNullInt64(&i)
	if actual.Int64 != i {
		t.Errorf("Expected %v but got %v", i, actual.Int64)
	}
	if actual.Valid != true {
		t.Errorf("Expected %v but got %v", true, actual.Valid)
	}
}

func Test_Uint64ToSqlNullInt64_NilPointer_ZeroAndInvalid(t *testing.T) {
	actual := Uint64ToSqlNullInt64(nil)
	if actual.Int64 != 0 {
		t.Errorf("Expected %v but got %v", 0, actual.Int64)
	}
	if actual.Valid != false {
		t.Errorf("Expected %v but got %v", false, actual.Valid)
	}
}

func Test_Uint64ToSqlNullInt64_NotNilPointer_Int64AndValid(t *testing.T) {
	i := uint64(10)
	actual := Uint64ToSqlNullInt64(&i)
	if actual.Int64 != int64(i) {
		t.Errorf("Expected %v but got %v", i, actual.Int64)
	}
	if actual.Valid != true {
		t.Errorf("Expected %v but got %v", true, actual.Valid)
	}
}

func Test_Uint32ToSqlNullInt32_NilPointer_ZeroAndInvalid(t *testing.T) {
	actual := Uint32ToSqlNullInt32(nil)
	if actual.Int32 != 0 {
		t.Errorf("Expected %v but got %v", 0, actual.Int32)
	}
	if actual.Valid != false {
		t.Errorf("Expected %v but got %v", false, actual.Valid)
	}
}

func Test_Uint32ToSqlNullInt32_NotNilPointer_Int32AndValid(t *testing.T) {
	i := uint32(10)
	actual := Uint32ToSqlNullInt32(&i)
	if actual.Int32 != int32(i) {
		t.Errorf("Expected %v but got %v", i, actual.Int32)
	}
	if actual.Valid != true {
		t.Errorf("Expected %v but got %v", true, actual.Valid)
	}
}

func Test_BoolToSqlNullBool_NilPointer_FalseAndInvalid(t *testing.T) {
	actual := BoolToSqlNullBool(nil)
	if actual.Bool != false {
		t.Errorf("Expected %v but got %v", false, actual.Bool)
	}
	if actual.Valid != false {
		t.Errorf("Expected %v but got %v", false, actual.Valid)
	}
}

func Test_BoolToSqlNullBool_NotNilPointer_BoolAndValid(t *testing.T) {
	b := true
	actual := BoolToSqlNullBool(&b)
	if actual.Bool != b {
		t.Errorf("Expected %v but got %v", b, actual.Bool)
	}
	if actual.Valid != true {
		t.Errorf("Expected %v but got %v", true, actual.Valid)
	}
}

func Test_Float64ToSqlNullFloat64_NilPointer_ZeroAndInvalid(t *testing.T) {
	actual := Float64ToSqlNullFloat64(nil)
	if actual.Float64 != 0 {
		t.Errorf("Expected %v but got %v", 0, actual.Float64)
	}
	if actual.Valid != false {
		t.Errorf("Expected %v but got %v", false, actual.Valid)
	}
}

func Test_Float64ToSqlNullFloat64_NotNilPointer_Float64AndValid(t *testing.T) {
	f := float64(10)
	actual := Float64ToSqlNullFloat64(&f)
	if actual.Float64 != f {
		t.Errorf("Expected %v but got %v", f, actual.Float64)
	}
	if actual.Valid != true {
		t.Errorf("Expected %v but got %v", true, actual.Valid)
	}
}

func Test_TimeToSqlNullTime_NilPointer_EmptyTimeAndInvalid(t *testing.T) {
	actual := TimeToSqlNullTime(nil)
	if !actual.Time.IsZero() {
		t.Errorf("Expected %v but got %v", "", actual.Time)
	}
	if actual.Valid != false {
		t.Errorf("Expected %v but got %v", false, actual.Valid)
	}
}

func Test_TimeToSqlNullTime_NotNilPointer_TimeAndValid(t *testing.T) {
	input := time.Now()
	actual := TimeToSqlNullTime(&input)
	if actual.Time != input {
		t.Errorf("Expected %v but got %v", t, actual.Time)
	}
	if actual.Valid != true {
		t.Errorf("Expected %v but got %v", true, actual.Valid)
	}
}

func Test_TimestampToSqlNullTime_NilPointer_EmptyTimeAndInvalid(t *testing.T) {
	actual := TimestampToSqlNullTime(nil)
	if !actual.Time.IsZero() {
		t.Errorf("Expected %v but got %v", "", actual.Time)
	}
	if actual.Valid != false {
		t.Errorf("Expected %v but got %v", false, actual.Valid)
	}
}

func Test_TimestampToSqlNullTime_NotNilPointer_TimeAndValid(t *testing.T) {
	input := timestamppb.Now()
	actual := TimestampToSqlNullTime(input)
	if actual.Time != input.AsTime() {
		t.Errorf("Expected %v but got %v", input.AsTime(), actual.Time)
	}
	if actual.Valid != true {
		t.Errorf("Expected %v but got %v", true, actual.Valid)
	}
}

func Test_SqlNullInt64ToInt64_InvalidNullInt64_Nil(t *testing.T) {
	i := sql.NullInt64{}
	actual := SqlNullInt64ToInt64(i)
	if actual != nil {
		t.Errorf("Expected %v but got %v", nil, actual)
	}
}

func Test_SqlNullInt64ToInt64_ValidNullInt64_Int64(t *testing.T) {
	i := sql.NullInt64{Int64: 10, Valid: true}
	actual := SqlNullInt64ToInt64(i)
	if *actual != i.Int64 {
		t.Errorf("Expected %v but got %v", i.Int64, *actual)
	}
}

func Test_SqlNullInt64ToUint64_ValidNullInt64_Uint64(t *testing.T) {
	i := sql.NullInt64{Int64: 10, Valid: true}
	actual := SqlNullInt64ToUint64(i)
	if *actual != uint64(i.Int64) {
		t.Errorf("Expected %v but got %v", uint64(i.Int64), *actual)
	}
}

func Test_SqlNullInt64ToUint64_InvalidNullInt64_Nil(t *testing.T) {
	i := sql.NullInt64{}
	actual := SqlNullInt64ToUint64(i)
	if actual != nil {
		t.Errorf("Expected %v but got %v", nil, actual)
	}
}

func Test_SqlNullStringToString_InvalidNullString_Nil(t *testing.T) {
	s := sql.NullString{}
	actual := SqlNullStringToString(s)
	if actual != nil {
		t.Errorf("Expected %v but got %v", nil, actual)
	}
}

func Test_SqlNullTimeToTimestamp_InvalidTime_Nil(t *testing.T) {
	input := sql.NullTime{}
	actual := SqlNullTimeToTimestamp(input)
	if actual != nil {
		t.Errorf("Expected %v but got %v", nil, actual)
	}
}

func Test_SqlNullTimeToTimestamp_ValidTime_Timestamp(t *testing.T) {
	input := sql.NullTime{Time: time.Now().UTC(), Valid: true}
	actual := SqlNullTimeToTimestamp(input)
	if actual.AsTime() != input.Time {
		t.Errorf("Expected %v but got %v", input.Time, actual.AsTime())
	}
}
