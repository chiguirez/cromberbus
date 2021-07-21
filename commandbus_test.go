package cromberbus_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	cb "github.com/chiguirez/cromberbus/v3"
	"github.com/stretchr/testify/require"

	"github.com/chiguirez/cromberbus/v3/commandhandler"
	"github.com/chiguirez/cromberbus/v3/middleware"
	"github.com/chiguirez/cromberbus/v3/middleware/multierror"
)

//nolint:funlen
func TestCommandBus(t *testing.T) {
	t.Run("Given a command a handler and a command bus", func(t *testing.T) {
		handler := commandhandler.New(func(ctx context.Context, a struct{}) error {
			return nil
		})
		commandBus := cb.New(handler)
		command := struct{}{}
		t.Run("When command goes into the bus", func(t *testing.T) {
			err := commandBus.Dispatch(context.Background(), command)
			t.Run("Then handler is executed", func(t *testing.T) {
				require.NoError(t, err)
			})
		})
	})

	t.Run("Given a command a handler and a command bus and a middleware", func(t *testing.T) {
		handlerErr := fmt.Errorf("handler")       //nolint:goerr113
		middlewareErr := fmt.Errorf("middleware") //nolint:goerr113
		handler := commandhandler.New(func(ctx context.Context, a struct{}) error {
			return handlerErr
		})
		commandBus := cb.New(handler)

		commandBus.Use(
			middleware.HandlerFunc(
				func(ctx context.Context, command middleware.Command, next middleware.NextFn) error {
					return multierror.New(
						next(ctx, command),
						middlewareErr,
					)
				},
			),
		)
		command := struct{}{}
		t.Run("When command goes into the bus", func(t *testing.T) {
			err := commandBus.Dispatch(context.Background(), command)
			t.Run("Then handler and middleware are executed", func(t *testing.T) {
				require.True(t, errors.Is(err, middlewareErr))
				require.True(t, errors.Is(err, handlerErr))
			})
		})
	})
}
