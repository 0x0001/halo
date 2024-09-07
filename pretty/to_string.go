package pretty

import (
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/jedib0t/go-pretty/v6/list"
	"github.com/jedib0t/go-pretty/v6/table"
)

func ToString(v interface{}, ignoreFieldNames ...string) string {
	t := reflect.TypeOf(v)
	switch t.Kind() {
	case reflect.Slice:
		return fromSlice(v, ignoreFieldNames...)
	case reflect.String:
		return reflect.ValueOf(v).String()
	case reflect.Map:
		return fromMap(v, ignoreFieldNames...)
	case reflect.Struct:
		return fromStruct(v, ignoreFieldNames...)
	case reflect.Ptr:
		v_ := reflect.ValueOf(v)
		return ToString(v_.Elem().Interface(), ignoreFieldNames...)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func fromSlice(v interface{}, ignoreFieldNames ...string) string {
	t := reflect.TypeOf(v)
	if t.Kind() != reflect.Slice {
		panic("not slice")
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
		lw := list.NewWriter()
		lw.SetStyle(list.StyleConnectedRounded)
		for i := 0; i < value.Len(); i++ {
			lw.AppendItem(ToString(value.Index(i).Interface()))
		}
		return lw.Render()
	}
}

func fromSliceStruct(v interface{}, ignoreFieldNames ...string) string {
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
				val = ToString(field.Interface())
			} else {
				val = field.Interface()
			}
			row = append(row, val)
		}
		w.AppendRow(table.Row(row))
		w.AppendSeparator()
	}
	return w.Render()
}

func fromSliceMap(v interface{}, ignoreFieldNames ...string) string {
	value := reflect.ValueOf(v)

	w := table.NewWriter()
	if value.IsNil() || value.IsZero() {
		return ""
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
		w.AppendSeparator()
	}
	return w.Render()
}

func fromMap(v interface{}, ignoreFieldNames ...string) string {
	t := reflect.TypeOf(v)
	if t.Kind() != reflect.Map {
		panic("not map")
	}
	value := reflect.ValueOf(v)

	w := table.NewWriter()
	if value.IsNil() || value.IsZero() {
		return ""
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
			row = append(row, ToString(v_.Interface()))
		}
	}
	w.AppendRow(table.Row(row))
	w.AppendSeparator()
	return w.Render()
}

func fromStruct(v interface{}, ignoreFieldNames ...string) string {
	value := reflect.ValueOf(v)
	t := value.Type()

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

	row := make([]any, 0, len(names))
	for _, fieldName := range names {
		field := value.FieldByName(fieldName.(string))
		row = append(row, ToString(field.Interface()))
	}
	w.AppendRow(table.Row(row))
	w.AppendSeparator()
	return w.Render()
}
