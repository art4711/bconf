// Copyright 2014 Thordur Bjornsson. All rights reserved.
// Use of this source code is governed by a ISC-style
// license that can be found in the LICENSE file.
package bconf

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Unmarshal the node, which must be a leaf node into v.
// If v is not a struct pointer, Unmarshal will panic.
func (bc Bconf) Unmarshal(v interface{}) (err error) {
	s := reflect.ValueOf(v).Elem()
	st := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		if !f.CanSet() {
			continue
		}

		name := st.Field(i).Name
		tag := st.Field(i).Tag.Get("bconf")
		if tag != "" {
			name = tag
		}
		name = strings.ToLower(name)
		value := bc.GetString(name)
		if value == "" {
			continue
		}

		switch f.Kind() {
		case reflect.String:
			f.SetString(value)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			ival, err := strconv.ParseInt(value, 0, f.Type().Bits())
			if err != nil {
				return err
			}
			f.SetInt(ival)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			uval, err := strconv.ParseUint(value, 0, f.Type().Bits())
			if err != nil {
				return err
			}
			f.SetUint(uval)
		case reflect.Float32, reflect.Float64:
			fval, err := strconv.ParseFloat(value, f.Type().Bits())
			if err != nil {
				return err
			}
			f.SetFloat(fval)
		default:
			return fmt.Errorf("invalid type %s\n", f.Kind())
		}
	}

	return nil
}
