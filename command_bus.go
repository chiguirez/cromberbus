package cromberbus

type Command interface{}

type CommandHandler interface {
	Handle(command Command) error
}

//TODO: what about error propagation and error handling
type CommandBus interface {
	Dispatch(command Command)
}
