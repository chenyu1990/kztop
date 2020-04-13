package kreedz

import (
	"fmt"
	"math"
)

func getInt64(a interface{}) int64 {
	switch a.(type) {
	case int:
		return int64(a.(int))
	case int64:
		return a.(int64)
	}
	return 0
}

func Inc(a interface{}, b interface{}) int64 {
	return getInt64(a) + getInt64(b)
}

func Mod(a interface{}, b int64) int64 {
	switch a.(type) {
	case int:
		return int64(a.(int)) % b
	case int64:
		return a.(int64) % b
	}
	return 0
}

func SubFloatRtnString(a float64, b float64) string {
	return fmt.Sprintf("%.2f", math.Abs(a - b))
}

func IsSlowly(a float64, b float64) bool {
	return a - b > 0
}