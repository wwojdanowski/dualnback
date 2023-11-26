package main

import (
	"fmt"
	"math/rand"
)

type Game struct {
	n           int
	boxQueue    []int
	letterQueue []int
	round       int
	maxRounds   int
	score       int
}

func (g *Game) nLastBox() int {
	return g.boxQueue[len(g.boxQueue)-1]
}

func (g *Game) nLastLetter() int {
	return g.letterQueue[len(g.letterQueue)-1]
}

func (g *Game) isReady() bool {
	return g.n == len(g.boxQueue)
}

func NewGame(n int, maxRounds int) *Game {
	g := Game{}
	g.n = n
	g.maxRounds = maxRounds
	g.boxQueue = make([]int, 0, n)
	g.letterQueue = make([]int, 0, n)
	return &g
}

func makeItem(box, letter int) Item {
	return Item{box: box, letter: letter}
}

func makeRandomItem() Item {
	return Item{box: makeBox(), letter: makeLetter()}
}

func makeBox() int {
	return rand.Intn(9)
}

func makeLetter() int {
	return rand.Intn(5)
}

func (g *Game) nextSequence(item Item) Item {
	b := item.box
	l := item.letter

	if g.isReady() {
		for i := g.n - 1; i > 0; i-- {
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
	box    int
	letter int
}

func (g *Game) scoreRound(hit bool) {
	if hit {
		g.score += 1
	}
	g.round += 1
}

func loop(g *Game, pulse <-chan struct{}, toggleBox <-chan struct{}, toggleLetter <-chan struct{}, feed <-chan Item) {
	boxPicked := false
	letterPicked := false
	lastItem := Item{}
	for {
		select {
		case <-pulse:
			if g.isReady() {
				if isCorrect(g, lastItem.box, lastItem.letter, boxPicked, letterPicked) {
					g.scoreRound(false)
				} else {
					g.scoreRound(true)
				}
				item := <-feed
				lastItem = g.nextSequence(item)
			} else {
				item := <-feed
				g.nextSequence(item)
			}
			// feed <- item

		case <-toggleBox:
			boxPicked = !boxPicked
		case <-toggleLetter:
			letterPicked = !letterPicked
		}
	}
	if g.round == g.maxRounds {

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

// func main2() {
// 	fmt.Println("Hello from Go Wasm!")
// 	js.Global().Get("document").Call("querySelector", "h1").Set("textContent", "Hello from Go Wasm!")
// 	select {}
// }
