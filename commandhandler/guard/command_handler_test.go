package guard_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	. "github.com/chiguirez/cromberbus/v3/commandhandler/guard"
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
	t.Run("Given an Interface that is a function but with more than one command struct", func(t *testing.T) {
		badCommandHandler := func(b struct{}) error { return nil }
		t.Run("When guarded then panic is returned", func(t *testing.T) {
			r.Panics(func() {
				CommandHandler(badCommandHandler)
			}, "handler should receive a context and a command struct only")
		})
	})
	t.Run("Given an Interface that is a function with and interface and a struct", func(t *testing.T) {
		badCommandHandler := func(a interface{}, b struct{}) error { return nil }
		t.Run("When guarded then panic is returned", func(t *testing.T) {
			r.Panics(func() {
				CommandHandler(badCommandHandler)
			}, "handler first argument should be a ctx")
		})
	})
	t.Run("Given an Interface that is a function with and context and an interface", func(t *testing.T) {
		badCommandHandler := func(a context.Context, b interface{}) error { return nil }
		t.Run("When guarded then panic is returned", func(t *testing.T) {
			r.Panics(func() {
				CommandHandler(badCommandHandler)
			}, "handler second argument must be a command struct")
		})
	})
	t.Run("Given an Interface that is a function with a context and an command struct and no return error", func(t *testing.T) {
		commandHandler := func(a context.Context, b struct{}) {}
		t.Run("When guarded then panic is returned", func(t *testing.T) {
			r.Panics(func() {
				CommandHandler(commandHandler)
			}, "handler should return error")
		})
	})
	t.Run("Given an Interface that is a function with a context and an command struct", func(t *testing.T) {
		commandHandler := func(a context.Context, b struct{}) error { return nil }
		t.Run("When guarded then no panic is returned", func(t *testing.T) {
			r.NotPanics(func() {
				CommandHandler(commandHandler)
			}, "should not return a panic")
		})
	})
}
