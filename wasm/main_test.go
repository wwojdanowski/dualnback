package main

import (
	"testing"
	"time"

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

	items := []Item{
		Item{0, 0},
		Item{1, 1},
		Item{2, 2},
		Item{3, 3},
		Item{4, 4},
		Item{5, 5},
		Item{6, 3},
		Item{4, 0},
		Item{8, 1},
		Item{8, 1},
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

	g.nextSequence(items[2])
	g.evalRound()

	assert.Equal(t, 2, g.score)
	g.nextSequence(items[3])
	g.evalRound()

	assert.Equal(t, 3, g.score)

	g.nextSequence(items[4])
	g.toggleBox()
	g.evalRound()

	assert.Equal(t, 3, g.score)

	g.nextSequence(items[5])
	g.toggleLetter()
	g.evalRound()

	assert.Equal(t, 3, g.score)

	g.nextSequence(items[6])
	g.toggleLetter()
	g.evalRound()

	assert.Equal(t, 4, g.score)

	g.nextSequence(items[7])
	g.toggleBox()
	g.evalRound()

	assert.Equal(t, 5, g.score)

	g.nextSequence(items[8])
	g.evalRound()
	assert.Equal(t, 6, g.score)

	assert.False(t, g.isDone())
	g.nextSequence(items[9])
	g.evalRound()
	assert.Equal(t, 7, g.score)

	assert.True(t, g.isDone())
}

func TestLoop(t *testing.T) {
	g := NewGame(3, 10)
	items := []Item{
		Item{0, 0},
		Item{1, 1},
		Item{2, 2},
		Item{3, 3},
		Item{4, 4},
		Item{5, 5},
		Item{6, 3},
		Item{4, 0},
		Item{8, 1},
		Item{8, 1},
		Item{0, 0},
		Item{1, 1},
		Item{2, 2},
		Item{3, 3},
	}

	ticker := time.NewTicker(50 * time.Millisecond)
	index := 0
	go func() {
		eval := false
		for {
			if index == len(items) {
				return
			}
			select {
			case <-ticker.C:
				if eval {
					if g.isReady() {
						g.evalRound()
					}
					eval = false
				} else {
					g.nextSequence(items[index])
					index++
					eval = true
				}
			}
		}
	}()

	time.Sleep(1500 * time.Millisecond)
	ticker.Stop()
	assert.Equal(t, g.round, 10)
	assert.Equal(t, g.score, 6)
	assert.True(t, g.isDone())

}
