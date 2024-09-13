package pretty

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/jedib0t/go-pretty/v6/list"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/stretchr/testify/suite"
)

type toStringTest struct {
	suite.Suite
}

func TestToString(t *testing.T) {
	suite.Run(t, new(toStringTest))
}

func (t *toStringTest) TestBasicTypes() {
	cases := map[string]any{
		"int8":    int8(1),
		"uint8":   uint8(1),
		"int16":   int16(1),
		"uint16":  uint16(1),
		"int32":   int32(1),
		"uint32":  uint32(1),
		"int64":   int64(1),
		"uint64":  uint64(1),
		"int":     int(1),
		"uint":    uint(1),
		"float32": float32(1),
		"float64": float64(1),
		"string":  "string",
		"nil":     nil,
	}

	for name, v := range cases {
		t.Run(name, func() {
			s, err := ToString(v)
			t.NoError(err)
			t.Equal(fmt.Sprintf("%v", v), s)
		})
	}
}

func (t *toStringTest) TestBasicPointerTypes() {
	i8 := int8(1)
	u8 := uint8(1)
	i16 := int16(1)
	u16 := uint16(1)
	i32 := int32(1)
	u32 := uint32(1)
	i64 := int64(1)
	u64 := uint64(1)
	i := int(1)
	u := uint(1)
	f32 := float32(1)
	f64 := float64(1)
	s := "string"

	cases := map[string]any{
		"int8":    &i8,
		"uint8":   &u8,
		"int16":   &i16,
		"uint16":  &u16,
		"int32":   &i32,
		"uint32":  &u32,
		"int64":   &i64,
		"uint64":  &u64,
		"int":     &i,
		"uint":    &u,
		"float32": &f32,
		"float64": &f64,
		"string":  &s,
		"nil":     nil,
	}

	for name, v := range cases {
		t.Run(name, func() {
			s, err := ToString(v)
			t.NoError(err)
			if v == nil {
				t.Equal("<nil>", s)
			} else {
				v_ := reflect.ValueOf(v).Elem().Interface()
				t.Equal(fmt.Sprintf("%v", v_), s)
			}
		})
	}

}

func (t *toStringTest) TestTime() {
	v := time.Now()
	s, err := ToString(v)
	t.NoError(err)
	t.Equal(v.Format(time.RFC3339), s)

	pv := &v
	s, err = ToString(pv)
	t.NoError(err)
	t.Equal(v.Format(time.RFC3339), s)
}

func (t *toStringTest) TestBasicSlice() {
	cases := map[string]any{
		"int8":    []int8{1, 2, 3},
		"uint8":   []uint8{1, 2, 3},
		"int16":   []int16{1, 2, 3},
		"uint16":  []uint16{1, 2, 3},
		"int32":   []int32{1, 2, 3},
		"uint32":  []uint32{1, 2, 3},
		"int64":   []int64{1, 2, 3},
		"uint64":  []uint64{1, 2, 3},
		"int":     []int{1, 2, 3},
		"uint":    []uint{1, 2, 3},
		"float32": []float32{1, 2, 3},
		"float64": []float64{1, 2, 3},
		"string":  []string{"a", "b", "c"},
		"nil":     nil,
	}

	for name, cases := range cases {
		t.Run(name, func() {
			s, err := ToString(cases)
			t.NoError(err)

			if cases == nil {
				t.Equal("<nil>", s)
			} else {
				w := list.NewWriter()
				w.SetStyle(list.StyleConnectedRounded)
				v := reflect.ValueOf(cases)
				for i := 0; i < v.Len(); i++ {
					w.AppendItem(fmt.Sprintf("%v", v.Index(i).Interface()))
				}

				t.Equal(w.Render(), s)
			}
		})
	}
}

func (t *toStringTest) TestStruct() {
	type S struct {
		A int
		B string
		C *string
		D time.Time
	}
	now := time.Now()

	v := S{1, "a", nil, now}
	s, err := ToString(v)
	t.NoError(err)

	w := table.NewWriter()
	w.AppendRow(table.Row{"A", 1})
	w.AppendRow(table.Row{"B", "a"})
	w.AppendRow(table.Row{"C", "<nil>"})
	w.AppendRow(table.Row{"D", now.Format(time.RFC3339)})

	t.Equal(w.Render(), s)

	pv := &v
	s, err = ToString(pv)
	t.NoError(err)
	t.Equal(w.Render(), s)
}

func (t *toStringTest) TestMap() {
	v := map[string]any{
		"a": 1,
		"d": "hello",
		"b": 2,
		"c": 3,
		"e": nil,
	}

	s, err := ToString(v)
	t.NoError(err)

	w := table.NewWriter()
	w.AppendRow(table.Row{"a", 1})
	w.AppendRow(table.Row{"b", 2})
	w.AppendRow(table.Row{"c", 3})
	w.AppendRow(table.Row{"d", "hello"})
	w.AppendRow(table.Row{"e", "<nil>"})

	t.Equal(w.Render(), s)

	pv := &v
	s, err = ToString(pv)
	t.NoError(err)
	t.Equal(w.Render(), s)
}

func (t *toStringTest) TestSliceStruct() {
	type S struct {
		A int
		B string
		C *string
	}

	v := []S{
		{1, "a", nil},
		{2, "b", nil},
		{3, "c", nil},
	}

	s, err := ToString(v)
	t.NoError(err)

	w := table.NewWriter()
	w.AppendHeader(table.Row{"A", "B", "C"})
	for _, s := range v {
		w.AppendRow(table.Row{s.A, s.B, "<nil>"})
	}

	t.Equal(w.Render(), s)

	pv := &v
	s, err = ToString(pv)
	t.NoError(err)
	t.Equal(w.Render(), s)
}

func (t *toStringTest) TestSlicePointerStruct() {
	type S struct {
		A int
		B string
		C *string
	}

	txt := "hello"
	v := []*S{
		{1, "a", nil},
		{2, "b", &txt},
		{3, "c", nil},
	}

	s, err := ToString(v)
	t.NoError(err)

	w := table.NewWriter()
	w.AppendHeader(table.Row{"A", "B", "C"})
	for _, s := range v {
		if s.C == nil {
			w.AppendRow(table.Row{s.A, s.B, "<nil>"})
		} else {
			w.AppendRow(table.Row{s.A, s.B, *s.C})
		}
	}

	t.Equal(w.Render(), s)

	pv := &v
	s, err = ToString(pv)
	t.NoError(err)
	t.Equal(w.Render(), s)
}

func (t *toStringTest) TestSliceMap() {
	v := []map[string]any{
		{"a": 1, "d": "hello"},
		{"b": 2, "c": 3},
		{"e": nil},
	}

	s, err := ToString(v)
	t.NoError(err)

	w := table.NewWriter()
	w.AppendHeader(table.Row{"a", "b", "c", "d", "e"})
	w.AppendRow(table.Row{1, "<nil>", "<nil>", "hello", "<nil>"})
	w.AppendRow(table.Row{"<nil>", 2, 3, "<nil>", "<nil>"})
	w.AppendRow(table.Row{"<nil>", "<nil>", "<nil>", "<nil>", "<nil>"})

	t.Equal(w.Render(), s)

	pv := &v
	s, err = ToString(pv)
	t.NoError(err)
	t.Equal(w.Render(), s)
}

func (t *toStringTest) TestSliceAny() {
	struct1 := struct {
		A int
		B string
	}{1, "a"}
	map1 := map[string]any{
		"a": 1,
		"d": "hello",
		"b": 2,
		"c": 3,
	}

	v := []any{
		struct1,
		&struct1,
		1,
		"hello",
		map1,
		&map1,
	}

	s, err := ToString(v)
	t.NoError(err)

	w := list.NewWriter()
	w.SetStyle(list.StyleConnectedRounded)
	{
		ww := table.NewWriter()
		ww.AppendRow(table.Row{"A", 1})
		ww.AppendRow(table.Row{"B", "a"})
		w.AppendItem(ww.Render())
		w.AppendItem(ww.Render())
	}
	w.AppendItem("1")
	w.AppendItem("hello")
	{
		ww := table.NewWriter()
		ww.AppendRow(table.Row{"a", 1})
		ww.AppendRow(table.Row{"b", 2})
		ww.AppendRow(table.Row{"c", 3})
		ww.AppendRow(table.Row{"d", "hello"})
		w.AppendItem(ww.Render())
		w.AppendItem(ww.Render())
	}

	t.Equal(w.Render(), s)

	pv := &v
	s, err = ToString(pv)
	t.NoError(err)
	t.Equal(w.Render(), s)
}

func (t *toStringTest) TestIntMap() {
	v := map[int]any{
		1: 1,
		2: "hello",
		3: 2,
		4: 3,
		5: nil,
	}

	s, err := ToString(v)
	t.NoError(err)

	w := table.NewWriter()
	w.AppendRow(table.Row{"1", 1})
	w.AppendRow(table.Row{"2", "hello"})
	w.AppendRow(table.Row{"3", 2})
	w.AppendRow(table.Row{"4", 3})
	w.AppendRow(table.Row{"5", "<nil>"})

	t.Equal(w.Render(), s)

	pv := &v
	s, err = ToString(pv)
	t.NoError(err)
	t.Equal(w.Render(), s)
}

func (t *toStringTest) TestJsonRawMessage() {
	var v json.RawMessage = []byte(`{"a":1,"b":"hello"}`)

	s, err := ToString(v)
	t.NoError(err)
	t.Equal(string(v), s)
}
