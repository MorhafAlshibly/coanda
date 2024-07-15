package conversion

import (
	"testing"
	"time"
)

func Test_TimeToTimestamppb_NilTime_NilTimestamp(t *testing.T) {
	actual := TimeToTimestamppb(nil)
	if actual != nil {
		t.Errorf("Expected %v but got %v", nil, actual)
	}
}

func Test_TimeToTimestamppb_NotNilTime_Timestamp(t *testing.T) {
	tm := time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)
	actual := TimeToTimestamppb(&tm)
	if !actual.AsTime().Equal(tm) {
		t.Errorf("Expected %v but got %v", tm, actual.AsTime())
	}
}

func Test_TimeToTimestamppb_ZeroTime_ZeroTimestamp(t *testing.T) {
	tm := time.Time{}
	actual := TimeToTimestamppb(&tm)
	if !actual.AsTime().Equal(tm) {
		t.Errorf("Expected %v but got %v", tm, actual.AsTime())
	}
}
