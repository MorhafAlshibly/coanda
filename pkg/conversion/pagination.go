package conversion

import "github.com/MorhafAlshibly/coanda/api"

func PaginationToLimitOffset(pagination *api.Pagination, defaultMax uint8, maxMax uint8) (uint64, uint64) {
	if pagination == nil {
		return uint64(defaultMax), 0
	}
	max := uint8(PointerToValue(pagination.Max, uint32(defaultMax)))
	if max > maxMax {
		max = maxMax
	}
	page := PointerToValue(pagination.Page, 1)
	offset := (page - 1) * uint64(max)
	return uint64(max), offset
}
