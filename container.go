package main

import (
	"errors"
)

type containerCallback func(c *container) interface{}

type container struct {
	services map[string]interface{}
}

var Container = &container{make(map[string]interface{})}

func (c *container) Bind(name string, callback containerCallback) {
	Container.services[name] = callback
}

func (c *container) Resolve(name string) interface{} {
	if _, ok := Container.services[name]; !ok {
		return errors.New("Service not found")
	}

	return Container.services[name].(containerCallback)(c)
}
