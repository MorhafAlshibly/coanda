package conversion

import "testing"

func Test_Enum_EmptyMap_DefaultValue(t *testing.T) {
	pbmap := map[string]int32{}
	dft := int32(8)
	actual := Enum("", pbmap, dft)
	if actual != dft {
		t.Errorf("Expected %v but got %v", dft, actual)
	}
}

func Test_Enum_ValueInMap_ValueInMap(t *testing.T) {
	pbmap := map[string]int32{"a": 1}
	dft := int32(8)
	actual := Enum("a", pbmap, dft)
	if actual != 1 {
		t.Errorf("Expected %v but got %v", 1, actual)
	}
}

func Test_Enum_ValueNotInMap_DefaultValue(t *testing.T) {
	pbmap := map[string]int32{"a": 1}
	dft := int32(8)
	actual := Enum("b", pbmap, dft)
	if actual != dft {
		t.Errorf("Expected %v but got %v", dft, actual)
	}
}

func Test_PointerToValue_NilPointer_DefaultValue(t *testing.T) {
	defaultValue := 8
	actual := PointerToValue(nil, defaultValue)
	if actual != defaultValue {
		t.Errorf("Expected %v but got %v", defaultValue, actual)
	}
}

func Test_PointerToValue_NotNilPointer_PointerValue(t *testing.T) {
	defaultValue := 8
	value := 10
	actual := PointerToValue(&value, defaultValue)
	if actual != value {
		t.Errorf("Expected %v but got %v", value, actual)
	}
}

func Test_ValueToPointer(t *testing.T) {
	value := 10
	actual := ValueToPointer(value)
	if *actual != value {
		t.Errorf("Expected %v but got %v", value, *actual)
	}
}
