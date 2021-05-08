package slice

import (
	"fmt"
	"reflect"
)

func Contains(target interface{}, list interface{}) (bool, error) {

	switch list.(type) {
	default:
		return false, fmt.Errorf("%v is an unsupported type", reflect.TypeOf(list))
	case []int:
		revert := list.([]int)
		for _, r := range revert {
			if target == r {
				return true, nil
			}
		}
		return false, nil

	case []int64:
		revert := list.([]uint64)
		for _, r := range revert {
			if target == r {
				return true, nil
			}
		}
		return false, nil

	case []string:
		revert := list.([]string)
		for _, r := range revert {
			if target == r {
				return true, nil
			}
		}
		return false, nil
	}

	return false, fmt.Errorf("processing failed")
}
