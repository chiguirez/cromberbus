package cromberbus

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type ACommandHandler struct {
	NumberOfHandleCalls int
}

func (h *ACommandHandler) Handle(command struct{}) error {
	h.NumberOfHandleCalls++
	return nil
}

type AMiddleware struct {
	NumberOfExecuteCalls int
}

func (m *AMiddleware) Execute(command Command, next CommandCallable) error {
	m.NumberOfExecuteCalls++

	return next(command)
}

func TestMapHandlerResolver_Resolve(t *testing.T) {
	isRequire := require.New(t)
	command := struct{}{}
	t.Run("Given a map handler resolver", func(t *testing.T) {
		sut := NewMapHandlerResolver()
		t.Run("When command handler is not found", func(t *testing.T) {
			handler, err := sut.Resolve(command)
			t.Run("Then an error is returned", func(t *testing.T) {
				isRequire.Equal(CommandHandler{}, handler)
				isRequire.Error(err)
			})
		})
		t.Run("When a command with its handler is added", func(t *testing.T) {
			handler := &ACommandHandler{}
			sut.AddHandler(handler.Handle)
			t.Run("Then the handler is resolved", func(t *testing.T) {
				resolvedHandler, err := sut.Resolve(command)
				isRequire.NoError(err)
				isRequire.Equal(
					NewCommandHandler(handler.Handle).CommandName(),
					resolvedHandler.CommandName(),
				)
			})
		})
	})
}

func TestCromberBus_Dispatch(t *testing.T) {
	isRequire := require.New(t)
	t.Run("Given a cromberbus command bus without middlewares", func(t *testing.T) {
		handler := ACommandHandler{}
		handlerResolver := CommandHandlerResolverMock{
			ResolveFunc: func(command Command) (CommandHandler, error) {
				return NewCommandHandler(handler.Handle), nil
			},
		}
		sut := NewCromberBus(&handlerResolver)
		t.Run("When a command is dispatched", func(t *testing.T) {
			command := struct{}{}
			_ = sut.Dispatch(command)
			t.Run("Then the resolved command handler handles the command", func(t *testing.T) {
				resolverHasBeenCalled := len(handlerResolver.ResolveCalls()) > 0
				isRequire.True(resolverHasBeenCalled)
				handlerHasBeenCalled := handler.NumberOfHandleCalls > 0
				isRequire.True(handlerHasBeenCalled)
			})
		})
	})
	t.Run("Given a cromberbus command bus with middlewares", func(t *testing.T) {
		aMiddleware := &AMiddleware{}
		anotherMiddleware := &AMiddleware{}
		handler := ACommandHandler{}
		handlerResolver := CommandHandlerResolverMock{
			ResolveFunc: func(command Command) (CommandHandler, error) {
				return NewCommandHandler(handler.Handle), nil
			},
		}
		sut := NewCromberBus(&handlerResolver, aMiddleware, anotherMiddleware)
		t.Run("When a command is dispatched", func(t *testing.T) {
			command := struct{}{}
			_ = sut.Dispatch(command)
			t.Run("Then the resolved command handler handles the command", func(t *testing.T) {
				resolverHasBeenCalled := len(handlerResolver.ResolveCalls()) > 0
				isRequire.True(resolverHasBeenCalled)
				handlerHasBeenCalled := handler.NumberOfHandleCalls > 0
				isRequire.True(handlerHasBeenCalled)
			})
			t.Run("And the middlewares are executed", func(t *testing.T) {
				isRequire.True(aMiddleware.NumberOfExecuteCalls > 0)
				isRequire.True(anotherMiddleware.NumberOfExecuteCalls > 0)
			})
		})
	})
}
