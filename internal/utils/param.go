package utils

import (
	"strconv"
)

func ParseUintParam(param string) uint {
	id, err := strconv.ParseUint(param, 10, 64)
	if err != nil {
		return 0
	}
	return uint(id)
}
