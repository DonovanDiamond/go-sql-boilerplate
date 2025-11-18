package boilerplate

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// These TypeSlices are used for postgresql arrays
//
// As PG uses {} for array syntax...

type StringSlice []string

func (s StringSlice) Value() (driver.Value, error) {
	if len(s) == 0 {
		return "{}", nil
	}
	return fmt.Sprintf(`{"%s"}`, strings.Join(s, `","`)), nil
}

func (s *StringSlice) Scan(src any) error {
	srcString := ""

	switch v := src.(type) {
	case string:
		srcString = v
	case []byte:
		srcString = string(v)
	default:
		return fmt.Errorf("unsupported type: %T for slice", src)
	}

	srcString = strings.Trim(srcString, "{}")
	if srcString == "" {
		return nil
	}
	*s = strings.Split(srcString, ",")

	// sqlite returns elements with quotes, so we need to get rid of them here
	for i := range *s {
		(*s)[i] = strings.Trim((*s)[i], `"`)
	}

	return nil
}

type IntSlice []int

func (s IntSlice) Value() (driver.Value, error) {
	if len(s) == 0 {
		return "{}", nil
	}
	v := ""
	for _, num := range s {
		v += fmt.Sprintf("%d,", num)
	}
	v = strings.TrimSuffix(v, ",")
	return fmt.Sprintf(`{%s}`, v), nil
}

func (s *IntSlice) Scan(src any) error {
	srcString := ""

	switch v := src.(type) {
	case string:
		srcString = v
	case []byte:
		srcString = string(v)
	default:
		return fmt.Errorf("unsupported type: %T for slice", src)
	}

	srcString = strings.Trim(srcString, "{}")
	if srcString == "" {
		return nil
	}
	nums := IntSlice{}
	for numString := range strings.SplitSeq(srcString, ",") {
		num, err := strconv.Atoi(numString)
		if err != nil {
			return err
		}
		nums = append(nums, num)
	}
	*s = nums

	return nil
}

type JsonObject map[string]any

func (obj JsonObject) Value() (driver.Value, error) {
	return json.Marshal(obj)
}

func (obj *JsonObject) Scan(src any) error {
	var srcBytes []byte

	switch v := src.(type) {
	case string:
		srcBytes = []byte(v)
	case []byte:
		srcBytes = v
	default:
		return fmt.Errorf("unsupported type: %T for json", src)
	}

	return json.Unmarshal(srcBytes, &obj)
}
