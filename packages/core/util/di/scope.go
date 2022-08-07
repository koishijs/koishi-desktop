package di

import "github.com/samber/do"

func Scope(i *do.Injector) (o *do.Injector) {
	o = do.New()
	services := i.ListProvidedServices()
	for _, service := range services {
		do.ProvideNamedValue(
			o,
			service,
			do.MustInvokeNamed[any](i, service),
		)
	}
	return
}
