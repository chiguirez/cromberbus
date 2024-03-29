package cromberbus

import (
	"context"

	"github.com/chiguirez/cromberbus/v3/commandhandler"
	"github.com/chiguirez/cromberbus/v3/maphandler"
	"github.com/chiguirez/cromberbus/v3/middleware"
)

type Command interface{}

type CommandBus struct {
	middlewareList middleware.List
	resolver       maphandler.Resolver
}

func (c *CommandBus) Use(middleware middleware.Handler) {
	c.middlewareList.Add(middleware)
}

func (c CommandBus) Dispatch(ctx context.Context, command Command) error {
	cHandler, err := c.resolver.Resolve(command)
	if err != nil {
		return err
	}

	return c.middlewareList.BuildWith(middleware.Wrap(cHandler)).Call(ctx, command)
}

func New(handlers ...commandhandler.Class) CommandBus {
	resolver := maphandler.NewResolver()

	for _, h := range handlers {
		resolver.AddHandler(h)
	}

	return CommandBus{
		middlewareList: make([]middleware.Handler, 0),
		resolver:       resolver,
	}
}
