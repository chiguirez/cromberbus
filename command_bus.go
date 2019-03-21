package cromberbus

type Command interface{}

type CommandHandler interface {
	Handle(command Command) error
}

type CommandBus interface {
	Dispatch(command Command) error
}
