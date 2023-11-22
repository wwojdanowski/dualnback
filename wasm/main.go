package main

import (
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

func makeBox() int {
	return rand.Intn(9)
}

func makeLetter() int {
	return rand.Intn(5)
}

func (g *Game) nextSequence() Item {
	b := makeBox()
	l := makeLetter()

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

	g.boxQueue[0] = b
	g.letterQueue[0] = l

	return Item{b, l}
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

func loop(g *Game, pulse <-chan struct{}, toggleBox <-chan struct{}, toggleLetter <-chan struct{}, feed chan<- Item) {
	boxPicked := false
	letterPicked := false
	presentedBox := -1
	presentedLetter := -1
	for {
		select {
		case <-pulse:
			if g.isReady() {
				if isCorrect(g, presentedBox, presentedLetter, boxPicked, letterPicked) {
					g.scoreRound(false)
				} else {
					g.scoreRound(true)
				}
			}
			item := g.nextSequence()
			feed <- item

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
	return match
}

// func main2() {
// 	fmt.Println("Hello from Go Wasm!")
// 	js.Global().Get("document").Call("querySelector", "h1").Set("textContent", "Hello from Go Wasm!")
// 	select {}
// }
