package network

import (
	"testing"
)

var ch = &channel{}

func TestChannel(t *testing.T) {

	t.Run("CreateChannel", func(t *testing.T) {
		ch.CreateChannel()
	})

	t.Run("JoinChannel", func(t *testing.T) {
		ch.JoinChannel()
	})

	t.Run("CreateChainCode", func(t *testing.T) {
		ch.CreateChainCode()
	})
}