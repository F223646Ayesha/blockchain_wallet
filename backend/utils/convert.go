package utils

import "strconv"

func FloatToString(v float64) string {
	return strconv.FormatFloat(v, 'f', -1, 64)
}

func IntToString(v int64) string {
	return strconv.FormatInt(v, 10)
}
