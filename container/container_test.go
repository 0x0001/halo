package container

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type containerTestSuit struct {
	suite.Suite
}

func TestContainer(t *testing.T) {
	suite.Run(t, new(containerTestSuit))
}

func (t *containerTestSuit) TestNew() {
	c := New()
	t.Assertions.NotNil(c, "New() should not return nil")
}

type bean1 struct{}
type bean2 struct {
	bean1 *bean1
}
type bean3 struct {
	bean1 *bean1
	bean2 *bean2
}

func (t *containerTestSuit) TestRegisterNoConstructor() {
	c := New()
	err := c.Register(BeanInfo{
		Name: "bean1",
	})
	t.Assertions.Error(err, "Register() should return error")
}

func (t *containerTestSuit) TestRegisterTwice() {
	c := New()
	err := c.Register(BeanInfo{
		Name: "bean1",
		Constructor: func(depends map[string]any, params map[string]any) (interface{}, error) {
			return &bean1{}, nil
		},
	})
	t.Assertions.NoError(err, "Register() should not return error")
	err = c.Register(BeanInfo{
		Name: "bean1",
		Constructor: func(depends map[string]any, params map[string]any) (interface{}, error) {
			return &bean1{}, nil
		},
	})
	t.Assertions.Error(err, "Register() should return error")
}

func (t *containerTestSuit) TestGet() {
	c := New()
	err := c.Register(BeanInfo{
		Name: "bean1",
		Constructor: func(depends map[string]any, params map[string]any) (interface{}, error) {
			return &bean1{}, nil
		},
	})
	t.Assertions.NoError(err, "Register() should not return error")
	bean, err := c.Get("bean1")
	t.Assertions.NoError(err, "Get() should not return error")
	_, ok := bean.(*bean1)
	t.Assertions.True(ok, "Get() should return bean1")
}

func (t *containerTestSuit) TestGetNotFound() {
	c := New()
	_, err := c.Get("bean1")
	t.Assertions.Error(err, "Get() should return error")
}

func (t *containerTestSuit) TestGetWithDependencies() {
	c := New()
	err := c.Register(BeanInfo{
		Name: "bean1",
		Constructor: func(depends map[string]any, params map[string]any) (interface{}, error) {
			return &bean1{}, nil
		},
	})
	t.Assertions.NoError(err, "Register() should not return error")
	err = c.Register(BeanInfo{
		Name:         "bean2",
		Dependencies: []string{"bean1"},
		Constructor: func(depends map[string]any, params map[string]any) (interface{}, error) {
			return &bean2{
				bean1: depends["bean1"].(*bean1),
			}, nil
		},
	})
	t.Assertions.NoError(err, "Register() should not return error")
	bean, err := c.Get("bean2")
	t.Assertions.NoError(err, "Get() should not return error")
	_, ok := bean.(*bean2)
	t.Assertions.True(ok, "Get() should return bean2")

	b2 := bean.(*bean2)
	t.Assertions.NotNil(b2.bean1, "bean1 should not be nil")

	t.Assertions.Equal(b2.bean1, c.beans["bean1"], "bean1 should be cached")
}

func (t *containerTestSuit) TestGetWithDependencies2() {
	c := New()
	err := c.Register(BeanInfo{
		Name: "bean1",
		Constructor: func(depends map[string]any, params map[string]any) (interface{}, error) {
			return &bean1{}, nil
		},
	})
	t.Assertions.NoError(err, "Register() should not return error")
	err = c.Register(BeanInfo{
		Name:         "bean2",
		Dependencies: []string{"bean1"},
		Constructor: func(depends map[string]any, params map[string]any) (interface{}, error) {
			return &bean2{
				bean1: depends["bean1"].(*bean1),
			}, nil
		},
	})
	t.Assertions.NoError(err, "Register() should not return error")
	err = c.Register(BeanInfo{
		Name:         "bean3",
		Dependencies: []string{"bean1", "bean2"},
		Constructor: func(depends map[string]any, params map[string]any) (interface{}, error) {
			return &bean3{
				bean1: depends["bean1"].(*bean1),
				bean2: depends["bean2"].(*bean2),
			}, nil
		},
	})
	t.Assertions.NoError(err, "Register() should not return error")
	bean, err := c.Get("bean3")
	t.Assertions.NoError(err, "Get() should not return error")
	_, ok := bean.(*bean3)
	t.Assertions.True(ok, "Get() should return bean3")

	b3 := bean.(*bean3)
	t.Assertions.NotNil(b3.bean1, "bean1 should not be nil")
	t.Assertions.NotNil(b3.bean2, "bean2 should not be nil")

	t.Assertions.Equal(b3.bean1, c.beans["bean1"], "bean1 should be cached")
	t.Assertions.Equal(b3.bean2, c.beans["bean2"], "bean2 should be cached")

}

type bean4 struct {
	bean5 *bean5
}
type bean5 struct {
	bean4 *bean4
}

func (t *containerTestSuit) TestGetWithCircularDependencies() {
	c := New()
	err := c.Register(BeanInfo{
		Name:         "bean4",
		Dependencies: []string{"bean5"},
		Constructor: func(depends map[string]any, params map[string]any) (interface{}, error) {
			return &bean4{
				bean5: depends["bean5"].(*bean5),
			}, nil
		},
	})
	t.Assertions.NoError(err, "Register() should not return error")
	err = c.Register(BeanInfo{
		Name:         "bean5",
		Dependencies: []string{"bean4"},
		Constructor: func(depends map[string]any, params map[string]any) (interface{}, error) {
			return &bean5{
				bean4: depends["bean4"].(*bean4),
			}, nil
		},
	})
	t.Assertions.NoError(err, "Register() should not return error")
	_, err = c.Get("bean4")
	t.Assertions.Error(err, "Get() should return error")
}

func (t *containerTestSuit) TestGetWithDependenciesNotFound() {
	c := New()
	err := c.Register(BeanInfo{
		Name:         "bean2",
		Dependencies: []string{"bean1"},
		Constructor: func(depends map[string]any, params map[string]any) (interface{}, error) {
			return &bean2{
				bean1: depends["bean1"].(*bean1),
			}, nil
		},
	})
	t.Assertions.NoError(err, "Register() should not return error")
	_, err = c.Get("bean2")
	t.Assertions.Error(err, "Get() should return error")
}

func (t *containerTestSuit) TestGetWithFunc() {
	c := New()
	err := c.Register(BeanInfo{
		Name: "bean1",
		Constructor: func(depends map[string]any, params map[string]any) (interface{}, error) {
			return &bean1{}, nil
		},
	})
	t.Assertions.NoError(err, "Register() should not return error")
	bean, err := Get[*bean1](c, "bean1")
	t.Assertions.NoError(err, "Get() should not return error")
	t.Assertions.NotNil(bean, "Get() should return bean1")
}

func (t *containerTestSuit) TestGetWithFuncNotFound() {
	c := New()
	_, err := Get[*bean1](c, "bean1")
	t.Assertions.Error(err, "Get() should return error")
}

func (t *containerTestSuit) TestGetWithCast() {
	c := New()
	err := c.Register(BeanInfo{
		Name: "bean1",
		Constructor: func(depends map[string]any, params map[string]any) (interface{}, error) {
			return &bean1{}, nil
		},
	})
	t.Assertions.NoError(err, "Register() should not return error")

	b, err := Get[*bean2](c, "bean1")
	t.Assertions.Error(err, "Get() should return error")
	t.Assertions.Nil(b, "Get() should return nil")
}
