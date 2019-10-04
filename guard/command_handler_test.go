package guard

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCommandHandler(t *testing.T) {
	r := require.New(t)
	t.Run("Given an Interface that is not a function", func(t *testing.T) {
		badCommandHandler := struct{}{}
		t.Run("When guarded then panic is returned", func(t *testing.T) {
			r.Panics(func() {
				CommandHandler(badCommandHandler)
			}, "handler should be a func")
		})
	})
	t.Run("Given an Interface that is a function but with more than one attribute", func(t *testing.T) {
		badCommandHandler := func(a struct{}, b struct{}) {}
		t.Run("When guarded then panic is returned", func(t *testing.T) {
			r.Panics(func() {
				CommandHandler(badCommandHandler)
			}, "handler should only receive one input")
		})
	})
	t.Run("Given an Interface that is a function with one attribute but its not a struct", func(t *testing.T) {
		badCommandHandler := func(a interface{}) {}
		t.Run("When guarded then panic is returned", func(t *testing.T) {
			r.Panics(func() {
				CommandHandler(badCommandHandler)
			}, "handler input must be a command struct")
		})
	})
}
