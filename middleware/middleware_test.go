package middleware_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/chiguirez/cromberbus/v3/commandhandler"
	"github.com/chiguirez/cromberbus/v3/middleware"
)

var ErrMiddleware = errors.New("error during middleware execution")

func TestMiddleware(t *testing.T) {
	t.Run("Given a List of Middleware", func(t *testing.T) {
		mw := make([]middleware.Handler, 0)
		list := middleware.List(mw)
		t.Run("When Add a 5 middleware as a pile", func(t *testing.T) {
			list.Add(loggerMiddleWare{t: t})
			list.Add(performanceMiddleWare{t: t})
			list.Add(
				middleware.Wrap(
					commandhandler.New(
						func(ctx context.Context, a struct{}) error {
							return fmt.Errorf("3rd Middleware : %w", ErrMiddleware)
						},
					),
				),
			)
			list.Add(
				middleware.Wrap(
					commandhandler.New(
						func(ctx context.Context, a struct{}) error {
							return fmt.Errorf("4th Middleware : %w", ErrMiddleware)
						},
					),
				),
			)
			list.Add(
				middleware.Wrap(
					commandhandler.New(
						func(ctx context.Context, a struct{}) error {
							return fmt.Errorf("5th Middleware : %w", ErrMiddleware)
						},
					),
				),
			)

			t.Run("Then errors all errors are dispatched ", func(t *testing.T) {
				cHandler := middleware.Wrap(
					commandhandler.New(
						func(ctx context.Context, a struct{}) error { return nil },
					),
				)
				err := list.BuildWith(cHandler).Call(context.Background(), struct{}{})
				require.Error(t, err)
				require.True(t, errors.Is(err, ErrMiddleware))
			})
		})
	})
}

type loggerMiddleWare struct {
	t *testing.T
}

func (l loggerMiddleWare) Handle(ctx context.Context, command middleware.Command, next middleware.NextFn) error {
	err := next(ctx, command)
	if err != nil {
		l.t.Logf("%+v", err)
	}

	return err
}

type performanceMiddleWare struct {
	t *testing.T
}

func (p performanceMiddleWare) Handle(ctx context.Context, command middleware.Command, next middleware.NextFn) error {
	now := time.Now()
	err := next(ctx, command)
	duration := time.Since(now)

	p.t.Logf("time elapsed: %s", duration.String())

	return err
}
