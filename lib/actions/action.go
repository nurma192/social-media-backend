package actions

import (
	"strconv"
)

func ParseToUint(value any) (uint, bool) {
	paramID, err := strconv.ParseUint(value.(string), 10, 64)
	if err != nil {
		return 0, false
	}

	return uint(paramID), true
}
