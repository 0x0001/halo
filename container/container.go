package container

import (
	"errors"
	"fmt"
)

type Container struct {
	infos map[string]BeanInfo
	beans map[string]any

	creating map[string]struct{}
}

func New() *Container {
	return &Container{
		beans:    make(map[string]any),
		infos:    make(map[string]BeanInfo),
		creating: make(map[string]struct{}),
	}
}

type Constructor func(depends map[string]any, params map[string]any) (interface{}, error)

type BeanInfo struct {
	Name         string
	Dependencies []string
	Params       map[string]any
	Constructor  Constructor
}

func (c *Container) Register(beanInfo BeanInfo) error {
	if _, ok := c.infos[beanInfo.Name]; ok {
		return errors.New("bean already exists: " + beanInfo.Name)
	}

	if beanInfo.Constructor == nil {
		return errors.New("constructor is required")
	}

	c.infos[beanInfo.Name] = beanInfo
	return nil
}

func (c *Container) Get(name string) (interface{}, error) {
	if bean, ok := c.beans[name]; ok {
		return bean, nil
	}
	beanInfo, ok := c.infos[name]
	if !ok {
		return nil, errors.New("bean not found: " + name)
	}

	// check circular dependency
	if _, ok := c.creating[name]; ok {
		return nil, errors.New("circular dependency: " + name)
	}
	c.creating[name] = struct{}{}
	defer delete(c.creating, name)

	depends := make(map[string]any)
	for _, dep := range beanInfo.Dependencies {
		depBean, err := c.Get(dep)
		if err != nil {
			return nil, err
		}
		depends[dep] = depBean
	}
	bean, err := beanInfo.Constructor(depends, beanInfo.Params)
	if err != nil {
		return nil, fmt.Errorf("create %s: %w", beanInfo.Name, err)
	}
	c.beans[name] = bean
	return bean, nil
}

func Get[T any](container *Container, name string) (T, error) {
	c, err := container.Get(name)
	if err != nil {
		var t T
		return t, err
	}

	if d, ok := c.(T); !ok {
		var t T
		return t, errors.New("type mismatch")
	} else {
		return d, nil
	}
}
