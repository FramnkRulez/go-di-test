package di

import (
	"errors"
	"fmt"
	"reflect"
)

// defines the di container to interface to struct and types to instances
type Container struct {
	registerMap map[reflect.Type]reflect.Type
	instanceMap map[reflect.Type]interface{}

	InitFunc string
}

// Register takes the interface type and the dependency type and maps them
func (c *Container) Register(i reflect.Type, t reflect.Type) error {
	if c.InitFunc == "" {
		c.InitFunc = "Init"
	}

	if c.registerMap == nil {
		c.registerMap = make(map[reflect.Type]reflect.Type)
	}

	_, key := c.registerMap[i]
	if key != false {
		return errors.New("interface already exists in container")
	}

	// _, exists := t.MethodByName(c.InitFunc)
	// if exists == false {
	// 	return errors.New(fmt.Sprintf("type being registered does not contain expected init func %s", c.InitFunc))
	// }

	c.registerMap[i] = t

	return nil
}

// Resolve returns a dependency by looking up its instance in the map and if it doesn't already exist,
// it will create the dependency (by resolving any required interfaces in the 'Init' method)
func (c *Container) Resolve(t reflect.Type) (interface{}, error) {
	dep := c.registerMap[t]

	if dep == nil {
		return nil, errors.New(fmt.Sprintf("failed to resolve dependency for interface type %T", dep))
	}

	if c.registerMap == nil {
		c.registerMap = make(map[reflect.Type]reflect.Type)
	}

	if c.instanceMap == nil {
		c.instanceMap = make(map[reflect.Type]interface{})
	}

	// check to see if we already created an instance, if so return it.
	inst, exists := c.instanceMap[t]

	// instance doesn't already exist, so we need to create it via reflection
	if exists == false {
		instPtr := reflect.New(dep)

		// find the 'init' method for this type
		init := instPtr.MethodByName(c.InitFunc)

		// create the inputs array for the method invoke
		inputs := make([]reflect.Value, init.Type().NumIn())

		// for every input into the 'init' method, resolve each dependency
		for i := 0; i < init.Type().NumIn(); i++ {
			input := init.Type().In(i)

			// call back into ourselves to resolve
			param, err := c.Resolve(input)

			// we failed to resolve a dependency, stop the resolve here
			if err != nil {
				return nil, err
			}

			inputs[i] = reflect.ValueOf(param)
		}

		// we were able to resolve all dependencies of the 'init' func, so invoke it now
		init.Call(inputs)
		inst = instPtr.Elem().Interface()

		// map this instance for futher dep resolution.
		c.instanceMap[t] = inst
	}

	return inst, nil
}
