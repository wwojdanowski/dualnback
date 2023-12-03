package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGame(t *testing.T) {
	g := NewGame(3, 10)

	g.nextSequence(makeRandomItem())
	g.nextSequence(makeRandomItem())
	g.nextSequence(makeRandomItem())
	g.nextSequence(makeRandomItem())

	assert.True(t, g.isReady(), "N-back is not ready!")
	assert.Len(t, g.boxQueue, 4)
	assert.Len(t, g.letterQueue, 4)
}

func TestGameNotReady(t *testing.T) {
	g := NewGame(3, 10)

	g.nextSequence(makeRandomItem())
	g.nextSequence(makeRandomItem())
	g.nextSequence(makeRandomItem())

	assert.False(t, g.isReady(), "N-back is ready!")
	assert.Len(t, g.boxQueue, 3)
	assert.Len(t, g.letterQueue, 3)

}

func TestNBackBeyondReady(t *testing.T) {
	g := NewGame(3, 10)

	g.nextSequence(makeRandomItem())
	g.nextSequence(makeRandomItem())
	g.nextSequence(makeRandomItem())
	g.nextSequence(makeRandomItem())
	g.nextSequence(makeRandomItem())
	g.nextSequence(makeRandomItem())

	assert.True(t, g.isReady(), "N-back is not ready!")
	assert.Len(t, g.boxQueue, 4)
	assert.Len(t, g.letterQueue, 4)
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
	feed <- makeRandomItem()
	pulse <- struct{}{}
	feed <- makeRandomItem()
	pulse <- struct{}{}
	feed <- makeRandomItem()
	pulse <- struct{}{}
	feed <- makeRandomItem()
	pulse <- struct{}{}
	feed <- makeRandomItem()
	toggleBox <- struct{}{}

	assert.True(t, g.isReady(), "N-back is not ready!")
	assert.Len(t, g.boxQueue, 4)
	assert.Len(t, g.letterQueue, 4)
}

func TestToggleCorrect(t *testing.T) {
	g := NewGame(3, 10)

	items := make([]Item, 10)

	for i := 0; i < len(items); i++ {
		items[i] = makeRandomItem()
	}

	g.nextSequence(items[0])
	g.nextSequence(items[1])
	g.nextSequence(items[2])
	g.nextSequence(items[0])
	assert.Equal(t, 0, g.score)

	g.toggleBox()
	g.toggleLetter()
	g.evalRound()

	assert.Equal(t, 1, g.score)

	g.nextSequence(items[1])
	g.toggleBox()
	g.toggleLetter()
	g.evalRound()

	assert.Equal(t, 2, g.score)
}
