package pretty

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/list"
	"github.com/jedib0t/go-pretty/v6/table"
)

func ToString(v interface{}, ignoreFieldNames ...string) (string, error) {
	t := reflect.TypeOf(v)
	switch t.Kind() {
	case reflect.Slice:
		return fromSlice(v, ignoreFieldNames...)
	case reflect.String:
		return reflect.ValueOf(v).String(), nil
	case reflect.Map:
		return fromMap(v, ignoreFieldNames...)
	case reflect.Struct:
		return fromStruct(v, ignoreFieldNames...)
	case reflect.Ptr:
		v_ := reflect.ValueOf(v)
		return ToString(v_.Elem().Interface(), ignoreFieldNames...)
	default:
		switch v.(type) {
		case time.Time:
			return v.(time.Time).Format(time.RFC3339), nil
		default:
			return fmt.Sprintf("%v", v), nil
		}
	}
}

func fromSlice(v interface{}, ignoreFieldNames ...string) (string, error) {
	t := reflect.TypeOf(v)
	if t.Kind() != reflect.Slice {
		return "", fmt.Errorf("not slice")
	}
	value := reflect.ValueOf(v)

	t = t.Elem()
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	switch t.Kind() {
	case reflect.Struct:
		return fromSliceStruct(v, ignoreFieldNames...)
	case reflect.Map:
		return fromSliceMap(v, ignoreFieldNames...)
	case reflect.Interface:
		fallthrough
	default:
		// []interface{}
		lw := list.NewWriter()
		lw.SetStyle(list.StyleConnectedRounded)
		for i := 0; i < value.Len(); i++ {
			v, err := ToString(value.Index(i).Interface())
			if err != nil {
				return "", err
			}
			lw.AppendItem(v)
		}
		return lw.Render(), nil
	}
}

func fromSliceStruct(v interface{}, ignoreFieldNames ...string) (string, error) {
	value := reflect.ValueOf(v)
	t := value.Type().Elem()

	// ptr := t.Kind() == reflect.Ptr
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
		value = value.Elem()
	}

	w := table.NewWriter()
	names := make([]any, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		name := t.Field(i).Name
		if slices.ContainsFunc(ignoreFieldNames, func(v string) bool { return strings.EqualFold(v, name) }) {
			continue
		}
		names = append(names, t.Field(i).Name)
	}
	w.AppendHeader(table.Row(names))

	for i := 0; i < value.Len(); i++ {
		row := make([]any, 0, t.NumField())
		for _, fieldName := range names {
			v := value.Index(i)
			field := v.FieldByName(fieldName.(string))
			var val any

			if field.Kind() == reflect.Slice {
				var err error
				val, err = ToString(field.Interface())
				if err != nil {
					return "", err
				}
			} else {
				val = field.Interface()
			}
			row = append(row, val)
		}
		w.AppendRow(table.Row(row))
	}
	return w.Render(), nil
}

func fromSliceMap(v interface{}, ignoreFieldNames ...string) (string, error) {
	value := reflect.ValueOf(v)

	w := table.NewWriter()
	if value.IsNil() || value.IsZero() {
		return "", nil
	}
	allKeys := make(map[string]struct{})
	for i := 0; i < value.Len(); i++ {
		for _, key := range value.Index(i).MapKeys() {
			name := key.String()
			if slices.ContainsFunc(ignoreFieldNames, func(v string) bool { return strings.EqualFold(v, name) }) {
				continue
			}
			allKeys[key.String()] = struct{}{}
		}
	}

	names := make([]any, 0, len(allKeys))
	for key := range allKeys {
		names = append(names, key)
	}
	slices.SortFunc(names, func(a, b any) int {
		return strings.Compare(a.(string), b.(string))
	})

	w.AppendHeader(table.Row(names))

	for i := 0; i < value.Len(); i++ {
		row := make([]any, 0, len(names))
		for _, key := range names {
			v_ := value.Index(i).MapIndex(reflect.ValueOf(key))
			if v_.Kind() == 0 {
				row = append(row, "")
			} else {
				row = append(row, v_.Interface())
			}
		}
		w.AppendRow(table.Row(row))
	}
	return w.Render(), nil
}

func fromMap(v interface{}, ignoreFieldNames ...string) (string, error) {
	t := reflect.TypeOf(v)
	if t.Kind() != reflect.Map {
		return "", fmt.Errorf("not map")
	}
	value := reflect.ValueOf(v)

	w := table.NewWriter()
	if value.IsNil() || value.IsZero() {
		return "", nil
	}
	allKeys := make(map[string]struct{})
	for _, key := range value.MapKeys() {
		name := key.String()
		if slices.ContainsFunc(ignoreFieldNames, func(v string) bool { return strings.EqualFold(v, name) }) {
			continue
		}
		allKeys[key.String()] = struct{}{}
	}

	names := make([]any, 0, len(allKeys))
	for key := range allKeys {
		names = append(names, key)
	}
	slices.SortFunc(names, func(a, b any) int {
		return strings.Compare(a.(string), b.(string))
	})

	w.AppendHeader(table.Row(names))

	row := make([]any, 0, len(names))
	for _, key := range names {
		v_ := value.MapIndex(reflect.ValueOf(key))
		if v_.Kind() == 0 {
			row = append(row, "")
		} else {
			if vv, err := ToString(v_.Interface()); err != nil {
				return "", err
			} else {
				row = append(row, vv)
			}
		}
	}
	w.AppendRow(table.Row(row))
	return w.Render(), nil
}

func fromStruct(v interface{}, ignoreFieldNames ...string) (string, error) {
	value := reflect.ValueOf(v)
	t := value.Type()
	if t == reflect.TypeOf(time.Time{}) {
		return value.Interface().(time.Time).Format(time.RFC3339), nil
	}

	w := table.NewWriter()
	names := make([]any, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		name := f.Name
		if slices.ContainsFunc(ignoreFieldNames, func(v string) bool { return strings.EqualFold(v, name) }) {
			continue
		}
		names = append(names, t.Field(i).Name)
	}
	w.AppendHeader(table.Row(names))

	row := make([]any, 0, len(names))
	for _, fieldName := range names {
		field := value.FieldByName(fieldName.(string))
		if field.CanInterface() {
			if vv, err := ToString(field.Interface()); err != nil {
				return "", err
			} else {
				row = append(row, vv)
			}
		} else {
			row = append(row, "")
		}
	}
	w.AppendRow(table.Row(row))
	return w.Render(), nil
}
