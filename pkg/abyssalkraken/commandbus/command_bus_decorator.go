package commandbus

type CommandBusDecorator interface {
	Decorate(command Command[any], execution func() (any, error)) func() (any, error)
}
