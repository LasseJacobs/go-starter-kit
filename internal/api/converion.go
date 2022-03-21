package api

import (
	"strconv"
)

func stringToInt(str string, min int64, max int64, onError int64) int64 {
	var intVal int64 = 25
	if str != "" {
		var err error
		intVal, err = strconv.ParseInt(str, 10, 0)
		if err != nil {
			return onError
		}
		if intVal < min {
			intVal = min
		}
		if intVal > max {
			intVal = max
		}
	}
	return intVal
}
