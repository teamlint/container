package di

import (
	"reflect"

	"github.com/teamlint/container/di/internal/reflection"
)

// newProviderInterface
func newProviderInterface(provider provider, as interface{}) *providerInterface {
	iface := reflection.InspectInterfacePtr(as)
	if !provider.Key().res.Implements(iface.Type) {
		panicf("%s not implement %s", provider.Key(), iface.Type)
	}
	return &providerInterface{
		res: key{
			name: provider.Key().name,
			res:  iface.Type,
			typ:  ptInterface,
		},
		provider: provider,
	}
}

// providerInterface
type providerInterface struct {
	res      key
	provider provider
}

func (i *providerInterface) Key() key {
	return i.res
}

func (i *providerInterface) ParameterList() parameterList {
	var plist parameterList
	plist = append(plist, parameter{
		name:     i.provider.Key().name,
		res:      i.provider.Key().res,
		optional: false,
		embed:    false,
	})
	return plist
}

func (i *providerInterface) Provide(parameters ...reflect.Value) (reflect.Value, error) {
	return parameters[0], nil
}
