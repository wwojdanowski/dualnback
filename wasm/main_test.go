package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGame(t *testing.T) {
	g := NewGame(3, 10)

	pulse := make(chan struct{})
	toggleBox := make(chan struct{})
	toggleLetter := make(chan struct{})
	feed := make(chan Item)

	go func() {
		loop(g, pulse, toggleBox, toggleLetter, feed)
	}()

	pulse <- struct{}{}
	pulse <- struct{}{}
	pulse <- struct{}{}

	assert.True(t, g.isReady(), "It is ready!")

}
