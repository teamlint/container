package di

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/teamlint/container/di/internal/reflection"
)

type ctorType int

const (
	ctorUnknown      ctorType = iota // unknown ctor signature
	ctorStd                          // (deps) (result)
	ctorError                        // (deps) (result, error)
	ctorCleanup                      // (deps) (result, cleanup)
	ctorCleanupError                 // (deps) (result, cleanup, error)
)

// newProviderConstructor
func newProviderConstructor(name string, ctor interface{}, history *provideHistory) *providerConstructor {
	if ctor == nil {
		panicf("The constructor must be a function like `func([dep1, dep2, ...]) (<result>, [cleanup, error])`, got `%s`", "nil")
	}
	if !reflection.IsFunc(ctor) {
		panicf("The constructor must be a function like `func([dep1, dep2, ...]) (<result>, [cleanup, error])`, got `%s`", reflect.ValueOf(ctor).Type())
	}
	fn := reflection.InspectFunction(ctor)
	ctorType := determineCtorType(fn)
	return &providerConstructor{
		name:     name,
		ctor:     fn,
		ctorType: ctorType,
		history:  history,
	}
}

// providerConstructor
type providerConstructor struct {
	name     string
	ctor     *reflection.Func
	ctorType ctorType
	clean    *reflection.Func
	history  *provideHistory
}

func (c providerConstructor) Key() key {
	return key{
		name: c.name,
		res:  c.ctor.Out(0),
		typ:  ptConstructor,
	}
}

func (c providerConstructor) ParameterList() parameterList {
	var plist parameterList
	for i := 0; i < c.ctor.NumIn(); i++ {
		ptype := c.ctor.In(i)
		var name string
		if ptype == parameterBagType {
			name = c.Key().String()
		}
		p := parameter{
			name:     name,
			res:      ptype,
			optional: false,
			embed:    isEmbedParameter(ptype),
		}
		plist = append(plist, p)
	}
	return plist
}

// Provide
func (c *providerConstructor) Provide(parameters ...reflect.Value) (reflect.Value, error) {
	defer func() {
		if c.history != nil {
			c.history.add(c.Key())
		}
	}()

	out := c.ctor.Call(parameters)
	switch c.ctorType {
	case ctorStd:
		return out[0], nil
	case ctorError:
		instance := out[0]
		err := out[1]
		if err.IsNil() {
			return instance, nil
		}
		return instance, err.Interface().(error)
	case ctorCleanup:
		c.setCleanup(out[1])
		return out[0], nil
	case ctorCleanupError:
		instance := out[0]
		cleanup := out[1]
		err := out[2]
		c.setCleanup(cleanup)
		if err.IsNil() {
			return instance, nil
		}
		return instance, err.Interface().(error)
	}
	return reflect.Value{}, errors.New("you found a bug, please create new issue for " +
		"this: https://github.com/defval/inject/issues/new")
}

func (c *providerConstructor) setCleanup(value reflect.Value) {
	c.clean = reflection.InspectFunction(value.Interface())
}

func (c *providerConstructor) Cleanup() {
	if c.clean != nil && c.clean.IsValid() {
		c.clean.Call([]reflect.Value{})
	}
}

// determineCtorType
func determineCtorType(fn *reflection.Func) ctorType {
	if fn.NumOut() == 1 {
		return ctorStd
	}
	if fn.NumOut() == 2 {
		if reflection.IsError(fn.Out(1)) {
			return ctorError
		}
		if reflection.IsCleanup(fn.Out(1)) {
			return ctorCleanup
		}
	}
	if fn.NumOut() == 3 && reflection.IsCleanup(fn.Out(1)) && reflection.IsError(fn.Out(2)) {
		return ctorCleanupError
	}
	panic(fmt.Sprintf("The constructor must be a function like `func([dep1, dep2, ...]) (<result>, [cleanup, error])`, got `%s`", fn.Name))
}
