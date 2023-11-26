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
	feed <- makeRandomItem()
	pulse <- struct{}{}
	feed <- makeRandomItem()
	pulse <- struct{}{}
	feed <- makeRandomItem()
	toggleBox <- struct{}{}

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
	feed <- makeRandomItem()
	pulse <- struct{}{}
	feed <- makeRandomItem()
	toggleBox <- struct{}{}

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
	feed <- makeRandomItem()
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

	items := make([]Item, 10)

	for i := 0; i < len(items); i++ {
		items[i] = makeRandomItem()
	}

	pulse <- struct{}{}
	feed <- items[0]
	pulse <- struct{}{}
	feed <- items[1]
	pulse <- struct{}{}
	feed <- items[2]
	pulse <- struct{}{}
	feed <- items[0]
	toggleBox <- struct{}{}
	toggleLetter <- struct{}{}

	assert.Equal(t, 0, g.score)
	pulse <- struct{}{}
	feed <- items[1]

	assert.Equal(t, 1, g.score)

	pulse <- struct{}{}
	feed <- items[2]

	assert.Equal(t, 2, g.score)

	pulse <- struct{}{}
	feed <- items[2]

}
