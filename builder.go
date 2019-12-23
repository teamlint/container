package container

import "github.com/teamlint/container/di"

var c *Container

// Instance current container instance
func Instance() *Container {
	return c
}

// Build creates a default container with provided options.
func Build(options ...Option) {
	c = &Container{
		container: di.New(),
	}
	// apply options.
	for _, opt := range options {
		opt.apply(c)
	}
	c.compile()
}

// Extract populates given target pointer with type instance provided in the container.
//
//   var server *http.Server
//   if err = container.Extract(&server); err != nil {
//     // extract failed
//   }
//
// If the target type does not exist in a container or instance type building failed, Extract() returns an error.
// Use ExtractOption for modifying the behavior of this function.
func Extract(target interface{}, options ...ExtractOption) (err error) {
	check()
	return c.Extract(target, options...)
}
func MustExtract(target interface{}, options ...ExtractOption) {
	check()
	c.MustExtract(target, options...)
}

// Invoke invokes custom function. Dependencies of function will be resolved via container.
func Invoke(fn interface{}) error {
	check()
	return c.Invoke(fn)
}
func MustInvoke(fn interface{}) {
	check()
	c.MustInvoke(fn)
}

// Cleanup cleanup container.
func Cleanup() {
	check()
	c.Cleanup()
}

func check() {
	if c == nil {
		panic("please building a container's instance")
	}
}
