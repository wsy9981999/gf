// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gvalid

import (
	"github.com/gogf/gf/internal/empty"
	"github.com/gogf/gf/util/gconv"
	"reflect"
	"strings"
)

// checkRequired checks `value` using required rules.
// It also supports require checks for `value` of type: slice, map.
func (v *Validator) checkRequired(value interface{}, ruleKey, rulePattern string, dataMap map[string]interface{}) bool {
	required := false
	switch ruleKey {
	// Required.
	case "required":
		required = true

	// Required unless all given field and its value are equal.
	// Example: required-if: id,1,age,18
	case "required-if":
		required = false
		array := strings.Split(rulePattern, ",")
		// It supports multiple field and value pairs.
		if len(array)%2 == 0 {
			for i := 0; i < len(array); {
				tk := array[i]
				tv := array[i+1]
				if v, ok := dataMap[tk]; ok {
					if strings.Compare(tv, gconv.String(v)) == 0 {
						required = true
						break
					}
				}
				i += 2
			}
		}

	// Required unless all given field and its value are not equal.
	// Example: required-unless: id,1,age,18
	case "required-unless":
		required = true
		array := strings.Split(rulePattern, ",")
		// It supports multiple field and value pairs.
		if len(array)%2 == 0 {
			for i := 0; i < len(array); {
				tk := array[i]
				tv := array[i+1]
				if v, ok := dataMap[tk]; ok {
					if strings.Compare(tv, gconv.String(v)) == 0 {
						required = false
						break
					}
				}
				i += 2
			}
		}

	// Required if any of given fields are not empty.
	// Example: required-with:id,name
	case "required-with":
		required = false
		array := strings.Split(rulePattern, ",")
		for i := 0; i < len(array); i++ {
			if !empty.IsEmpty(dataMap[array[i]]) {
				required = true
				break
			}
		}

	// Required if all of given fields are not empty.
	// Example: required-with:id,name
	case "required-with-all":
		required = true
		array := strings.Split(rulePattern, ",")
		for i := 0; i < len(array); i++ {
			if empty.IsEmpty(dataMap[array[i]]) {
				required = false
				break
			}
		}

	// Required if any of given fields are empty.
	// Example: required-with:id,name
	case "required-without":
		required = false
		array := strings.Split(rulePattern, ",")
		for i := 0; i < len(array); i++ {
			if empty.IsEmpty(dataMap[array[i]]) {
				required = true
				break
			}
		}

	// Required if all of given fields are empty.
	// Example: required-with:id,name
	case "required-without-all":
		required = true
		array := strings.Split(rulePattern, ",")
		for i := 0; i < len(array); i++ {
			if !empty.IsEmpty(dataMap[array[i]]) {
				required = false
				break
			}
		}
	}
	if required {
		reflectValue := reflect.ValueOf(value)
		for reflectValue.Kind() == reflect.Ptr {
			reflectValue = reflectValue.Elem()
		}
		switch reflectValue.Kind() {
		case reflect.String, reflect.Map, reflect.Array, reflect.Slice:
			return reflectValue.Len() != 0
		}
		return gconv.String(value) != ""
	} else {
		return true
	}
}
