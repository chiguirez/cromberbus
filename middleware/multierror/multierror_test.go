package multierror_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/chiguirez/cromberbus/v3/middleware/multierror"
)

var (
	ErrCustom     = errors.New("custom error")
	ErrAnother    = errors.New("another error")
	ErrNotThisOne = errors.New("not this one")
)

func TestMultiErr_Is(t *testing.T) {
	t.Run("Given multiple errors", func(t *testing.T) {
		t.Run("When merged into multiErr struct", func(t *testing.T) {
			err := multierror.New(ErrCustom, ErrAnother)
			t.Run("Then we can still check for any error type", func(t *testing.T) {
				require.True(t, errors.Is(err, ErrCustom))
				require.True(t, errors.Is(err, ErrAnother))
				require.True(t, !errors.Is(err, ErrNotThisOne))
			})
		})
	})
}
