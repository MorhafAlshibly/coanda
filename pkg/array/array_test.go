package array

import (
	"reflect"
	"testing"
)

func Test_RemoveDuplicates_EmptyList_EmptyList(t *testing.T) {
	list := []int{}
	expected := []int{}
	actual := RemoveDuplicates(list)
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v but got %v", expected, actual)
	}
}

func Test_RemoveDuplicates_ListWithDuplicates_ListWithoutDuplicates(t *testing.T) {
	list := []int{1, 2, 3, 1, 2, 3}
	expected := []int{1, 2, 3}
	actual := RemoveDuplicates(list)
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v but got %v", expected, actual)
	}
}

func Test_RemoveDuplicates_ListWithoutDuplicates_ListWithoutDuplicates(t *testing.T) {
	list := []int{1, 2, 3}
	expected := []int{1, 2, 3}
	actual := RemoveDuplicates(list)
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v but got %v", expected, actual)
	}
}

func Test_Remove_EmptyList_EmptyList(t *testing.T) {
	list := []int{}
	expected := []int{}
	actual := Remove(list, 1)
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v but got %v", expected, actual)
	}
}

func Test_Remove_ListWithoutItem_ListWithoutItem(t *testing.T) {
	list := []int{1, 2, 3}
	expected := []int{1, 2, 3}
	actual := Remove(list, 4)
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v but got %v", expected, actual)
	}
}

func Test_Remove_ListWithItem_ListWithoutItem(t *testing.T) {
	list := []int{1, 2, 3}
	expected := []int{1, 3}
	actual := Remove(list, 2)
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v but got %v", expected, actual)
	}
}

func Test_Contains_EmptyList_False(t *testing.T) {
	list := []int{}
	expected := false
	actual := Contains(list, 1)
	if expected != actual {
		t.Errorf("Expected %v but got %v", expected, actual)
	}
}

func Test_Contains_ListWithoutItem_False(t *testing.T) {
	list := []int{1, 2, 3}
	expected := false
	actual := Contains(list, 4)
	if expected != actual {
		t.Errorf("Expected %v but got %v", expected, actual)
	}
}

func Test_Contains_ListWithItem_True(t *testing.T) {
	list := []int{1, 2, 3}
	expected := true
	actual := Contains(list, 2)
	if expected != actual {
		t.Errorf("Expected %v but got %v", expected, actual)
	}
}
