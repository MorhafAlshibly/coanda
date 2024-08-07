package conversion

import (
	"testing"

	"github.com/MorhafAlshibly/coanda/api"
)

func Test_PaginationToLimitOffset_NilPagination_DefaultValues(t *testing.T) {
	defaultMax := uint8(8)
	maxMax := uint8(10)
	actualMax, actualOffset := PaginationToLimitOffset(nil, defaultMax, maxMax)
	if actualMax != uint64(defaultMax) {
		t.Errorf("Expected %v but got %v", defaultMax, actualMax)
	}
	if actualOffset != 0 {
		t.Errorf("Expected %v but got %v", 0, actualOffset)
	}
}

func Test_PaginationToLimitOffset_MaxGreaterThanMaxMax_MaxMax(t *testing.T) {
	defaultMax := uint8(8)
	maxMax := uint8(10)
	max := uint32(11)
	pagination := &api.Pagination{Max: &max}
	actualMax, actualOffset := PaginationToLimitOffset(pagination, defaultMax, maxMax)
	if actualMax != uint64(maxMax) {
		t.Errorf("Expected %v but got %v", maxMax, actualMax)
	}
	if actualOffset != 0 {
		t.Errorf("Expected %v but got %v", 0, actualOffset)
	}
}

func Test_PaginationToLimitOffset_MaxLessThanMaxMax_PaginationMax(t *testing.T) {
	defaultMax := uint8(8)
	maxMax := uint8(10)
	max := uint32(9)
	pagination := &api.Pagination{Max: &max}
	actualMax, actualOffset := PaginationToLimitOffset(pagination, defaultMax, maxMax)
	if actualMax != 9 {
		t.Errorf("Expected %v but got %v", max, actualMax)
	}
	if actualOffset != 0 {
		t.Errorf("Expected %v but got %v", 0, actualOffset)
	}
}

func Test_PaginationToLimitOffset_PageNotSet_MaxAndZeroOffset(t *testing.T) {
	defaultMax := uint8(8)
	maxMax := uint8(10)
	max := uint32(9)
	pagination := &api.Pagination{Max: &max}
	actualMax, actualOffset := PaginationToLimitOffset(pagination, defaultMax, maxMax)
	if actualMax != 9 {
		t.Errorf("Expected %v but got %v", max, actualMax)
	}
	if actualOffset != 0 {
		t.Errorf("Expected %v but got %v", 0, actualOffset)
	}
}

func Test_PaginationToLimitOffset_PageSet_MaxAndOffset(t *testing.T) {
	defaultMax := uint8(8)
	maxMax := uint8(10)
	max := uint32(9)
	page := uint64(2)
	pagination := &api.Pagination{Max: &max, Page: &page}
	actualMax, actualOffset := PaginationToLimitOffset(pagination, defaultMax, maxMax)
	if actualMax != 9 {
		t.Errorf("Expected %v but got %v", max, actualMax)
	}
	if actualOffset != 9 {
		t.Errorf("Expected %v but got %v", max, actualOffset)
	}
}
