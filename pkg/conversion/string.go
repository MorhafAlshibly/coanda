package conversion

import "fmt"

// Convert a uint64 array to comma separated string
func Uint64ArrayToCommaSeparatedString(arr []uint64) string {
	if len(arr) == 0 {
		return ""
	}
	var str string
	for i, v := range arr {
		if i == 0 {
			str = str + fmt.Sprintf("%d", v)
		} else {
			str = str + "," + fmt.Sprintf("%d", v)
		}
	}
	return str
}
