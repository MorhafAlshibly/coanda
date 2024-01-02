package validation

func ValidateMaxPageLength(max *uint32, defaultMaxPageLength uint8, maxMaxPageLength uint8) uint8 {
	if max == nil {
		return defaultMaxPageLength
	}
	if *max > uint32(maxMaxPageLength) {
		return maxMaxPageLength
	}
	return uint8(*max)
}
