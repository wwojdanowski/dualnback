package game

import (
	"fmt"
	"math/rand"
)

type GameObserver interface {
	NewSequence(*Game, Item)
	PauseForDecision(*Game)
	EvalRound(*Game)
	RoundFinished(*Game)
	StateProcessed(*Game)
	ToggleBox(*Game)
	ToggleLetter(*Game)
}

type Game struct {
	N              int
	Round          int
	MaxRounds      int
	LastResult     Result
	Score          int
	boxQueue       []int
	letterQueue    []int
	boxSelected    bool
	letterSelected bool
}

func (g *Game) nLastBox() int {
	return g.boxQueue[len(g.boxQueue)-1]
}

func (g *Game) nLastLetter() int {
	return g.letterQueue[len(g.letterQueue)-1]
}

func (g *Game) firstBox() int {
	return g.boxQueue[0]
}

func (g *Game) firstLetter() int {
	return g.letterQueue[0]
}

func (g *Game) IsReady() bool {
	return g.N+1 == len(g.boxQueue)
}

func (g *Game) isCompleted() bool {
	return g.MaxRounds == g.Round
}

func NewGame(n int, maxRounds int) *Game {
	g := Game{}
	g.N = n
	g.MaxRounds = maxRounds
	g.boxQueue = make([]int, 0, n+1)
	g.letterQueue = make([]int, 0, n+1)
	return &g
}

func makeItem(box, letter int) Item {
	return Item{Box: box, Letter: letter}
}

func MakeRandomItem() Item {
	return Item{Box: makeBox(), Letter: makeLetter()}
}

func makeBox() int {
	return rand.Intn(9)
}

func makeLetter() int {
	return rand.Intn(5)
}

func (g *Game) IsBoxToggled() bool {
	return g.boxSelected
}

func (g *Game) IsLetterToggled() bool {
	return g.letterSelected
}

func (g *Game) ToggleBox() {
	g.boxSelected = !g.boxSelected
}

func (g *Game) ToggleLetter() {
	g.letterSelected = !g.letterSelected
}

type Result struct {
	Box    bool
	Letter bool
}

func (g *Game) EvalRound() {
	result := Result{true, true}
	score := false
	g.Round += 1

	if g.nLastBox() == g.firstBox() {
		score = g.boxSelected
	} else {
		score = !g.boxSelected
	}

	if score {
		if g.nLastLetter() == g.firstLetter() {
			score = g.letterSelected
		} else {
			score = !g.letterSelected
		}
	} else {
		result.Box = false
	}

	if score {
		g.Score += 1
	} else {
		result.Letter = false
	}

	g.LastResult = result
	g.resetToggles()
}

func (g *Game) IsDone() bool {
	return g.Round == g.MaxRounds
}

func (g *Game) resetToggles() {
	g.boxSelected = false
	g.letterSelected = false
}

func (g *Game) NextSequence(item Item) Item {
	b := item.Box
	l := item.Letter

	if g.IsReady() {
		for i := g.N; i > 0; i-- {
			g.boxQueue[i] = g.boxQueue[i-1]
			g.letterQueue[i] = g.letterQueue[i-1]
		}
	} else {
		g.boxQueue = g.boxQueue[:len(g.boxQueue)+1]
		g.letterQueue = g.letterQueue[:len(g.letterQueue)+1]

		for i := len(g.boxQueue) - 1; i > 0; i-- {
			g.boxQueue[i] = g.boxQueue[i-1]
			g.letterQueue[i] = g.letterQueue[i-1]
		}
	}

	ret := Item{g.boxQueue[0], g.letterQueue[0]}
	g.boxQueue[0] = b
	g.letterQueue[0] = l
	return ret
}

type Item struct {
	Box    int
	Letter int
}

func (g *Game) scoreRound(hit bool) {
	if hit {
		g.Score += 1
	}
	g.Round += 1
}

func loop(g *Game, pulse <-chan struct{}, toggleBox <-chan struct{}, toggleLetter <-chan struct{}, feed <-chan Item) {
	boxPicked := false
	letterPicked := false
	lastItem := Item{}
	for {
		select {
		case <-pulse:
			if g.IsReady() {
				if isCorrect(g, lastItem.Box, lastItem.Letter, boxPicked, letterPicked) {
					g.scoreRound(false)
				} else {
					g.scoreRound(true)
				}
				item := <-feed
				lastItem = g.NextSequence(item)
			} else {
				item := <-feed
				g.NextSequence(item)
			}
			// feed <- item

		case <-toggleBox:
			boxPicked = !boxPicked
		case <-toggleLetter:
			letterPicked = !letterPicked
		}
	}
	if g.Round == g.MaxRounds {

	}
}

func isCorrect(g *Game, presentedBox, presentedLetter int, boxPicked, letterPicked bool) bool {
	match := true
	if presentedBox == g.nLastBox() {
		if !boxPicked {
			match = false
		}
	} else {
		if boxPicked {
			match = false
		}
	}
	if match {
		if presentedLetter == g.nLastLetter() {
			if !letterPicked {
				match = false
			}
		} else {
			if letterPicked {
				match = false
			}
		}
	}

	fmt.Println(presentedBox, presentedLetter, boxPicked, letterPicked, match)

	return match
}
