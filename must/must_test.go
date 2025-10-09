package must

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	type testCase struct {
		input  any
		result any
	}

	tests := []testCase{
		{"string", "string"},
		{1, 1},
		{1.0, 1.0},
		{true, true},
	}

	for _, tt := range tests {
		v := Get(tt.input, nil)
		assert.Equal(t, tt.result, v)
	}

	func() {
		defer func() {
			if r := recover(); r != nil {
				assert.Equal(t, "has error", r.(error).Error())
			}
		}()
		Get(1, errors.New("has error"))
		t.Fail()
	}()
}
