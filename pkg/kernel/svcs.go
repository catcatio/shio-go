package kernel

type Initializer struct {
	Name     string
	Register func(config *ServiceOptions)
	OnServe  func(config *ServiceOptions)
}
