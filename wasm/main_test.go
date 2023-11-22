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
	<-feed
	pulse <- struct{}{}
	<-feed
	pulse <- struct{}{}
	<-feed

	assert.True(t, g.isReady(), "N-back is not ready!")
	assert.Len(t, g.boxQueue, 3)
	assert.Len(t, g.letterQueue, 3)
}

func TestGameNotReady(t *testing.T) {
	g := NewGame(3, 10)

	pulse := make(chan struct{})
	toggleBox := make(chan struct{})
	toggleLetter := make(chan struct{})
	feed := make(chan Item)

	go func() {
		loop(g, pulse, toggleBox, toggleLetter, feed)
	}()

	pulse <- struct{}{}
	<-feed
	pulse <- struct{}{}
	<-feed

	assert.False(t, g.isReady(), "N-back is ready!")
	assert.Len(t, g.boxQueue, 2)
	assert.Len(t, g.letterQueue, 2)

}

func TestNBackBeyondReady(t *testing.T) {
	g := NewGame(3, 10)

	pulse := make(chan struct{})
	toggleBox := make(chan struct{})
	toggleLetter := make(chan struct{})
	feed := make(chan Item)

	go func() {
		loop(g, pulse, toggleBox, toggleLetter, feed)
	}()

	pulse <- struct{}{}
	<-feed
	pulse <- struct{}{}
	<-feed
	pulse <- struct{}{}
	<-feed
	pulse <- struct{}{}
	<-feed
	pulse <- struct{}{}
	<-feed
	pulse <- struct{}{}
	<-feed

	assert.True(t, g.isReady(), "N-back is not ready!")
	assert.Len(t, g.boxQueue, 3)
	assert.Len(t, g.letterQueue, 3)
}

func TestSelect(t *testing.T) {
	g := NewGame(3, 10)

	pulse := make(chan struct{})
	toggleBox := make(chan struct{})
	toggleLetter := make(chan struct{})
	feed := make(chan Item)

	go func() {
		loop(g, pulse, toggleBox, toggleLetter, feed)
	}()

	pulse <- struct{}{}
	<-feed
	pulse <- struct{}{}
	<-feed
	pulse <- struct{}{}
	<-feed
	pulse <- struct{}{}
	<-feed
	pulse <- struct{}{}
	<-feed
	pulse <- struct{}{}
	<-feed

	assert.True(t, g.isReady(), "N-back is not ready!")
	assert.Len(t, g.boxQueue, 3)
	assert.Len(t, g.letterQueue, 3)
}

func TestToggleCorrect(t *testing.T) {
	g := NewGame(3, 10)

	pulse := make(chan struct{})
	toggleBox := make(chan struct{})
	toggleLetter := make(chan struct{})
	feed := make(chan Item)

	go func() {
		loop(g, pulse, toggleBox, toggleLetter, feed)
	}()

	pulse <- struct{}{}
	<-feed
	pulse <- struct{}{}
	<-feed
	pulse <- struct{}{}
	<-feed
	pulse <- struct{}{}

	
}
