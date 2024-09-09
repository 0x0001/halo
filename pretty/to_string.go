package pretty

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/list"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

func ToString(v interface{}, ignoreFieldNames ...string) (string, error) {
	if t, ok := v.(time.Time); ok {
		return t.Format(time.RFC3339), nil
	}
	if v == nil {
		return "<nil>", nil
	}
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
		if v_.IsValid() && !v_.IsNil() {
			return ToString(v_.Elem().Interface(), ignoreFieldNames...)
		}
		return fmt.Sprintf("%v", v), nil
	default:
		return fmt.Sprintf("%v", v), nil
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

	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	w := table.NewWriter()
	names := make([]any, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		name := t.Field(i).Name
		if containsIgnoreCase(ignoreFieldNames, name) {
			continue
		}
		names = append(names, t.Field(i).Name)
	}
	w.AppendHeader(names)

	for i := 0; i < value.Len(); i++ {
		row := make([]any, 0, t.NumField())
		for _, fieldName := range names {
			vv := value.Index(i)
			for vv.Kind() == reflect.Ptr {
				vv = vv.Elem()
			}
			field := vv.FieldByName(fieldName.(string))
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
			rv, err := ToString(val, ignoreFieldNames...)
			if err != nil {
				return "", err
			}
			row = append(row, rv)
		}
		w.AppendRow(row)
	}
	return w.Render(), nil
}

func fromSliceMap(v interface{}, ignoreFieldNames ...string) (string, error) {
	value := reflect.ValueOf(v)

	w := table.NewWriter()
	if !value.IsValid() || value.Len() == 0 {
		return "", nil
	}

	allKeys := make(map[string]struct{})
	for i := 0; i < value.Len(); i++ {
		for _, key := range value.Index(i).MapKeys() {
			name := key.String()
			if containsIgnoreCase(ignoreFieldNames, name) {
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

	w.AppendHeader(names)

	for i := 0; i < value.Len(); i++ {
		row := make([]any, 0, len(names))
		for _, key := range names {
			v_ := value.Index(i).MapIndex(reflect.ValueOf(key))
			if !v_.IsValid() {
				row = append(row, "<nil>")
			} else {
				row = append(row, v_.Interface())
			}
		}
		w.AppendRow(row)
	}
	return w.Render(), nil
}

func fromMap(v interface{}, ignoreFieldNames ...string) (string, error) {
	t := reflect.TypeOf(v)
	if t.Kind() != reflect.Map {
		return "", fmt.Errorf("not map")
	}
	value := reflect.ValueOf(v)

	if value.IsNil() || value.IsZero() {
		return "", nil
	}
	allKeys := make(map[interface{}]struct{})
	for _, key := range value.MapKeys() {
		name := fmt.Sprintf("%v", key.Interface())
		if containsIgnoreCase(ignoreFieldNames, name) {
			continue
		}
		allKeys[key.Interface()] = struct{}{}
	}

	names := make([]any, 0, len(allKeys))
	for key := range allKeys {
		names = append(names, key)
	}
	slices.SortFunc(names, func(a, b any) int {
		keyA := fmt.Sprintf("%v", a)
		keyB := fmt.Sprintf("%v", b)
		return strings.Compare(keyA, keyB)
	})

	w := table.NewWriter()
	w.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, Align: text.AlignRight},
	})

	for _, key := range names {
		v_ := value.MapIndex(reflect.ValueOf(key))
		if v_.IsValid() {
			if vv, err := ToString(v_.Interface()); err != nil {
				return "", err
			} else {
				w.AppendRow(table.Row{key, vv})
			}
		}
	}
	return w.Render(), nil
}

func fromStruct(v interface{}, ignoreFieldNames ...string) (string, error) {
	value := reflect.ValueOf(v)
	t := value.Type()

	names := make([]any, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		name := f.Name
		if containsIgnoreCase(ignoreFieldNames, name) {
			continue
		}
		names = append(names, t.Field(i).Name)
	}

	w := table.NewWriter()
	w.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, Align: text.AlignRight},
	})

	for _, fieldName := range names {
		field := value.FieldByName(fieldName.(string))
		if field.CanInterface() {
			if vv, err := ToString(field.Interface()); err != nil {
				return "", err
			} else {
				w.AppendRow(table.Row{fieldName, vv})
			}
		}
	}
	return w.Render(), nil
}

func containsIgnoreCase(slice []string, str string) bool {
	for _, s := range slice {
		if strings.EqualFold(s, str) {
			return true
		}
	}
	return false
}
