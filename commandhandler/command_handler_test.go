package commandhandler_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/chiguirez/cromberbus/v3/commandhandler"
)

var ErrCustom = errors.New("custom Error")

func TestCommandHandler(t *testing.T) {
	type customClass struct{}

	called := make(chan struct{}, 1)

	t.Run("Given a valid Handler function", func(t *testing.T) {
		handlerFn := func(ctx context.Context, a customClass) error {
			called <- struct{}{}

			return ErrCustom
		}
		handler := commandhandler.New(handlerFn)
		t.Run("When Call function is invoked", func(t *testing.T) {
			err := handler.Call(context.Background(), customClass{})
			t.Run("Then Handler gets called", func(t *testing.T) {
				<-called
				require.True(t, errors.Is(err, ErrCustom))
				require.Equal(t, handler.CommandName(), "customClass")
			})
		})
	})
}
