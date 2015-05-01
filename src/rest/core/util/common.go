package util

import (
	"fmt"
	"strconv"
)

func GetIntValue(v interface{}) *int {
	switch t := v.(type) {
	case float32:
		x := int(t)
		return &x
	case float64:
		x := int(t)
		return &x
	case int:
		return &t
	case nil:
		return nil
	case string:
		if i, err := strconv.Atoi(t); err != nil {
			return nil
		} else {
			return &i
		}
	default:
		fmt.Print("Failed to get int default")
		return nil
	}
}
