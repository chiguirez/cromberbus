package middleware

import (
	"context"

	"github.com/chiguirez/cromberbus/v3/commandhandler"
	"github.com/chiguirez/cromberbus/v3/middleware/multierror"
)

type Command interface{}

type Class struct {
	handler Handler
	next    NextFn
}

func (m Class) Call(ctx context.Context, command Command) error {
	return m.handler.Handle(ctx, command, m.next)
}

type NextFn func(ctx context.Context, cmd Command) error

type Handler interface {
	Handle(ctx context.Context, command Command, next NextFn) error
}

type HandlerFunc func(ctx context.Context, command Command, next NextFn) error

func (h HandlerFunc) Handle(ctx context.Context, command Command, next NextFn) error {
	return h(ctx, command, next)
}

func Wrap(handler commandhandler.Class) Handler {
	return HandlerFunc(func(ctx context.Context, command Command, next NextFn) error {
		return multierror.New(
			handler.Call(ctx, command),
			next(ctx, command),
		).NilOrError()
	})
}

type List []Handler

func (l *List) Add(h Handler) {
	*l = append(*l, h)
}

func (l List) BuildWith(handler Handler) Class {
	clone := append(l[:0:0], l...)
	clone.Add(handler)

	return build(clone)
}

func build(handlers []Handler) Class {
	var next Class

	switch {
	case len(handlers) == 0:
		return void()
	case len(handlers) > 1:
		next = build(handlers[1:])
	default:
		next = void()
	}

	return Class{handlers[0], next.Call}
}

func void() Class {
	return Class{
		handler: HandlerFunc(func(ctx context.Context, command Command, next NextFn) error { return nil }),
		next:    nil,
	}
}
