package typer

import (
	"github.com/stretchr/testify/require"
	"testing"
)

type ACommand struct {}

func TestIdentify(t *testing.T) {
	assertThat := require.New(t)

	t.Run("Given a ACommand variable", func(t *testing.T) {
		command := ACommand{}
		t.Run("When its type is identified", func(t *testing.T) {
			result := Identify(command)
			t.Run("Then a 'ACommand' name is returned", func(t *testing.T) {
				assertThat.Equal("ACommand", result)
			})
		})
	})
	t.Run("Given a pointer to an ACommand variable", func(t *testing.T) {
		command := new(ACommand)
		t.Run("When its type is identified", func(t *testing.T) {
			result := Identify(command)
			t.Run("Then a 'ACommand' name is returned", func(t *testing.T) {
				assertThat.Equal("ACommand", result)
			})
		})
	})
}
