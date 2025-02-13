package di

import "reflect"

// provider lookup sequence
var providerLookupSequence = []providerType{ptConstructor, ptInterface, ptGroup, ptEmbedParameter}

// providerType
type providerType int

const (
	ptUnknown providerType = iota
	ptConstructor
	ptInterface
	ptGroup
	ptEmbedParameter
)

// provideHistory
type provideHistory struct {
	items []key
}

// add
func (h *provideHistory) add(k key) {
	h.items = append(h.items, k)
}

// provider
type provider interface {
	// The identity of result type.
	Key() key
	// ParameterList returns array of dependencies.
	ParameterList() parameterList
	// Provide provides value from provided parameters.
	Provide(values ...reflect.Value) (reflect.Value, error)
}

// cleanup
type cleanup interface {
	Cleanup()
}
