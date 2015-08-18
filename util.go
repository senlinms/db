// Copyright 2014 The zhgo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

import (
	"errors"
	"log"
	"reflect"
)

// Get scan variables
func scanVariables(ptr interface{}, columnsLen int, isRows bool) (reflect.Kind, interface{}, []interface{}, error) {
	typ := reflect.TypeOf(ptr)
	if typ.Kind() != reflect.Ptr {
		return 0, nil, nil, errors.New("ptr is not a pointer")
	}

	elemTyp := typ.Elem()

	if isRows { // Rows
		if elemTyp.Kind() != reflect.Slice {
			return 0, nil, nil, errors.New("ptr is not point a slice")
		}

		elemTyp = elemTyp.Elem()
	}

	elemKind := elemTyp.Kind()

	// element(value) is point to row
	scan := make([]interface{}, columnsLen)

	if elemKind == reflect.Struct {
		if columnsLen != elemTyp.NumField() {
			return 0, nil, nil, errors.New("columnsLen is not equal elemTyp.NumField()")
		}

		row := reflect.New(elemTyp) // Data
		for i := 0; i < columnsLen; i++ {
			f := elemTyp.Field(i)
			if !f.Anonymous { // && f.Tag.Get("field") != ""
				scan[i] = row.Elem().FieldByIndex([]int{i}).Addr().Interface()
			}
		}

		return elemKind, row.Interface(), scan, nil
	}

	if elemKind == reflect.Map || elemKind == reflect.Slice {
		row := make([]interface{}, columnsLen) // Data
		for i := 0; i < columnsLen; i++ {
			scan[i] = &row[i]
		}

		return elemKind, &row, scan, nil
	}

	return 0, nil, nil, errors.New("ptr is not a point struct, map or slice")
}

// Type assertions
func typeAssertion(v interface{}) interface{} {
	switch v.(type) {
	case []byte:
		return v.([]byte)
	case []rune:
		return v.([]rune)
	case bool:
		return v.(bool)
	case float64:
		return v.(float64)
	case int64:
		return v.(int64)
	case nil:
		return nil
	case string:
		return v.(string)
	default:
		log.Printf("Unexpected type %#v\n", v)
		return ""
	}
}

// Table alias
func tableAlias(alias []string) string {
	if len(alias) > 0 {
		return alias[0]
	}
	return ""
}

// Reflect struct, construct Field slice
func tableFields(entity interface{}) (string, []string, []string, map[string]string) {
	typ := reflect.Indirect(reflect.ValueOf(entity)).Type()
	primary := ""
	fields := make([]string, 0)
	allFields := make([]string, 0)
	jsonMap := make(map[string]string)

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		var fd string
		if field.Tag.Get("field") != "" {
			fd = field.Tag.Get("field")
		} else {
			fd = field.Name
		}

		var jn string
		if field.Tag.Get("json") != "" {
			jn = field.Tag.Get("json")
		} else {
			jn = field.Name
		}

		//!field.Anonymous
		if field.Tag.Get("pk") == "true" {
			primary = fd
		} else {
			fields = append(fields, fd)
		}

		allFields = append(allFields, fd)
		jsonMap[jn] = fd
	}

	return primary, fields, allFields, jsonMap
}
