package conversion

func Enum[T ~string, PB ~int32](val T, pbmap map[string]int32, dft PB) PB {
	v, ok := pbmap[string(val)]
	if !ok {
		return dft
	}
	return PB(v)
}

// Pointer to Value
func PointerToValue[T any](p *T, defaultValue T) T {
	if p == nil {
		return defaultValue
	}
	return *p
}

// Value to Pointer
func ValueToPointer[T any](v T) *T {
	return &v
}

func PointerBoolToValue(p *bool) bool {
	if p == nil {
		return false
	}
	return *p
}
