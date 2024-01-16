package wasm

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGame(t *testing.T) {
	g := NewGame(3, 10)

	g.NextSequence(MakeRandomItem())
	g.NextSequence(MakeRandomItem())
	g.NextSequence(MakeRandomItem())
	g.NextSequence(MakeRandomItem())

	assert.True(t, g.IsReady(), "N-back is not ready!")
	assert.Len(t, g.boxQueue, 4)
	assert.Len(t, g.letterQueue, 4)
}

func TestGameNotReady(t *testing.T) {
	g := NewGame(3, 10)

	g.NextSequence(MakeRandomItem())
	g.NextSequence(MakeRandomItem())
	g.NextSequence(MakeRandomItem())

	assert.False(t, g.IsReady(), "N-back is ready!")
	assert.Len(t, g.boxQueue, 3)
	assert.Len(t, g.letterQueue, 3)

}

func TestNBackBeyondReady(t *testing.T) {
	g := NewGame(3, 10)

	g.NextSequence(MakeRandomItem())
	g.NextSequence(MakeRandomItem())
	g.NextSequence(MakeRandomItem())
	g.NextSequence(MakeRandomItem())
	g.NextSequence(MakeRandomItem())
	g.NextSequence(MakeRandomItem())

	assert.True(t, g.IsReady(), "N-back is not ready!")
	assert.Len(t, g.boxQueue, 4)
	assert.Len(t, g.letterQueue, 4)
}

func TestSelect(t *testing.T) {
	g := NewGame(3, 10)

	pulse := make(chan struct{})
	ToggleBox := make(chan struct{})
	ToggleLetter := make(chan struct{})
	feed := make(chan Item)

	go func() {
		loop(g, pulse, ToggleBox, ToggleLetter, feed)
	}()

	pulse <- struct{}{}
	feed <- MakeRandomItem()
	pulse <- struct{}{}
	feed <- MakeRandomItem()
	pulse <- struct{}{}
	feed <- MakeRandomItem()
	pulse <- struct{}{}
	feed <- MakeRandomItem()
	pulse <- struct{}{}
	feed <- MakeRandomItem()
	ToggleBox <- struct{}{}

	assert.True(t, g.IsReady(), "N-back is not ready!")
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

	g.NextSequence(items[0])
	g.NextSequence(items[1])
	g.NextSequence(items[2])
	g.NextSequence(items[0])
	assert.Equal(t, 0, g.Score)

	g.ToggleBox()
	g.ToggleLetter()
	g.EvalRound()

	assert.Equal(t, 1, g.Score)

	g.NextSequence(items[1])
	g.ToggleBox()
	g.ToggleLetter()
	g.EvalRound()

	assert.Equal(t, 2, g.Score)

	g.NextSequence(items[2])
	g.EvalRound()

	assert.Equal(t, 2, g.Score)
	g.NextSequence(items[3])
	g.EvalRound()

	assert.Equal(t, 3, g.Score)

	g.NextSequence(items[4])
	g.ToggleBox()
	g.EvalRound()

	assert.Equal(t, 3, g.Score)

	g.NextSequence(items[5])
	g.ToggleLetter()
	g.EvalRound()

	assert.Equal(t, 3, g.Score)

	g.NextSequence(items[6])
	g.ToggleLetter()
	g.EvalRound()

	assert.Equal(t, 4, g.Score)

	g.NextSequence(items[7])
	g.ToggleBox()
	g.EvalRound()

	assert.Equal(t, 5, g.Score)

	g.NextSequence(items[8])
	g.EvalRound()
	assert.Equal(t, 6, g.Score)

	assert.False(t, g.IsDone())
	g.NextSequence(items[9])
	g.EvalRound()
	assert.Equal(t, 7, g.Score)

	assert.True(t, g.IsDone())
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
					if g.IsReady() {
						g.EvalRound()
					}
					eval = false
				} else {
					g.NextSequence(items[index])
					index++
					eval = true
				}
			}
		}
	}()

	time.Sleep(1500 * time.Millisecond)
	ticker.Stop()
	assert.Equal(t, g.Round, 10)
	assert.Equal(t, g.Score, 6)
	assert.True(t, g.IsDone())
}
