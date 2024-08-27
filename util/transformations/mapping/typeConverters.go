// Package mapping is for best guess data transforms, where
// some sort of mapping has been provided to help fill in the blanks.
package mapping

import (
	"strconv"
)

// AnyToInt converts any number to an integer, including floats with truncations
// Non numbers are returned as 0
func AnyToInt(in any) int {
	switch dint := in.(type) {
	// @TODO convert the floats to integers
	case int:
		return int(dint)
	case int64:
		return int(dint)
	case int16:
		return int(dint)
	case int32:
		return int(dint)
	case int8:
		return int(dint)
	case float64:
		return int(dint)
	case float32:
		return int(dint)
	case string:
		if inInt, err := strconv.Atoi(dint); err == nil {
			return inInt
		} else {
			return 0
		}
	default:

		return 0
	}

}

// AnyToFloat64 converts an any to a float. Non numbers are returned as 0
func AnyToFloat64(in any) float64 {

	switch float := in.(type) {
	case int:
		return float64(float)
	case int64:
		return float64(float)
	case int16:
		return float64(float)
	case int32:
		return float64(float)
	case int8:
		return float64(float)
	case float64:
		return float64(float)
	case float32:
		return float64(float)
	default:
		// fmt.Println(float)
		return 0
	}
}
