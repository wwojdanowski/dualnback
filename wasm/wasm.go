package main

import (
	"dualnback/game"
	"fmt"
	"syscall/js"
	"time"
)

type JSGameObserver struct {
	done   chan bool
	canvas js.Value
}

/*

const canvas = document.getElementById("myCanvas");
const ctx = canvas.getContext("2d");


ctx.font = "30px Arial";
ctx.fillText("Hello World", 10, 50);

*/

func drawGridWithItem(canvas js.Value, newItem game.Item) interface{} {
	context := canvas.Call("getContext", "2d")

	boxRow := newItem.Box / 3
	boxCol := newItem.Box % 3

	cellSize := 100
	cellMargin := 20
	letters := []string{"A", "B", "C", "D", "E"}

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			x := j * cellSize
			y := i * cellSize
			if i == boxRow && j == boxCol {
				context.Set("fillStyle", "green")
				context.Call("fillRect", x+cellMargin, y+cellMargin,
					cellSize-cellMargin, cellSize-cellMargin)
				font := fmt.Sprintf("%dpx Arial", int(float64(cellSize)*0.6))
				context.Set("font", font)
				context.Set("strokeStyle", "black")
				context.Call("strokeText", letters[newItem.Letter], x+cellMargin+cellSize/4, y+cellMargin+cellSize/2)
				context.Set("fillStyle", "black")
			} else {
				context.Call("fillRect", x+cellMargin, y+cellMargin,
					cellSize-cellMargin, cellSize-cellMargin)
			}
		}
	}

	return js.Undefined()
}

func drawGrid(canvas js.Value) interface{} {
	context := canvas.Call("getContext", "2d")

	cellSize := 100
	cellMargin := 20
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			x := j * cellSize
			y := i * cellSize
			context.Call("fillRect", x+cellMargin, y+cellMargin,
				cellSize-cellMargin, cellSize-cellMargin)
		}
	}

	return js.Undefined()
}

func (o *JSGameObserver) NewSequence(g *game.Game, item game.Item) {
	drawGridWithItem(o.canvas, item)
	if g.IsReady() {
		drawButtons(o.canvas, "ready")
	} else {
		drawButtons(o.canvas, "gray")
	}
}

func (o *JSGameObserver) PauseForDecision(g *game.Game) {
	drawGrid(o.canvas)
}

func drawPlaceButton(canvas js.Value, state string) {
	context := canvas.Call("getContext", "2d")

	var style string
	switch state {
	case "ready":
		style = "blue"
	case "correct":
		style = "green"
	case "wrong":
		style = "red"
	case "toggle":
		style = "yellow"
	default:
		style = "gray"
	}

	context.Set("fillStyle", style)
	context.Call("fillRect", 10, 400, 150, 50)

	context.Set("font", "30px serif")
	context.Call("strokeText", "PLACE", 20, 430)
	context.Set("fillStyle", "black")
}

func drawTextButton(canvas js.Value, state string) {
	context := canvas.Call("getContext", "2d")

	var style string
	switch state {
	case "ready":
		style = "blue"
	case "correct":
		style = "green"
	case "wrong":
		style = "red"
	case "toggle":
		style = "yellow"
	default:
		style = "gray"
	}

	context.Set("fillStyle", style)
	context.Call("fillRect", 170, 400, 150, 50)

	context.Set("font", "30px serif")
	context.Call("strokeText", "LETTER", 180, 430)
	context.Set("fillStyle", "black")
}

func printScoreBoard(canvas js.Value, g *game.Game) {
	context := canvas.Call("getContext", "2d")
	context.Set("font", "30px serif")
	scoreBar := fmt.Sprintf("N: %d | score: %d | rounds: %d/%d", g.N, g.Score, g.Round, g.MaxRounds)
	context.Set("fillStyle", "white")
	context.Call("fillRect", 0, 470, 500, 50)
	context.Set("fillStyle", "black")
	context.Call("strokeText", scoreBar, 20, 500)
}

func (o *JSGameObserver) EvalRound(g *game.Game) {

	if g.IsReady() {
		if g.LastResult.Box {
			drawPlaceButton(o.canvas, "correct")
		} else {
			drawPlaceButton(o.canvas, "wrong")
		}

		if g.LastResult.Letter {
			drawTextButton(o.canvas, "correct")
		} else {
			drawTextButton(o.canvas, "wrong")
		}
	}

	printScoreBoard(o.canvas, g)
}

func (o *JSGameObserver) RoundFinished(g *game.Game) {
	drawButtons(o.canvas, "ready")
}

func (o *JSGameObserver) StateProcessed(g *game.Game) {
}

func (o *JSGameObserver) ToggleBox(g *game.Game) {
	if g.IsBoxToggled() {
		drawPlaceButton(o.canvas, "toggle")
	} else {
		drawPlaceButton(o.canvas, "ready")
	}
}

func (o *JSGameObserver) ToggleLetter(g *game.Game) {
	if g.IsLetterToggled() {
		drawTextButton(o.canvas, "toggle")
	} else {
		drawTextButton(o.canvas, "ready")
	}
}

func drawTestGrid(canvas js.Value) interface{} {
	newItem := game.Item{1, 1}
	return drawGridWithItem(canvas, newItem)
}

func run(this js.Value, p []js.Value) interface{} {
	g := game.NewGame(2, 20)
	ticker := time.NewTicker(500 * time.Millisecond)
	toggleBox := make(chan struct{})
	toggleLetter := make(chan struct{})
	done := make(chan bool)

	document := js.Global().Get("document")
	canvas := document.Call("createElement", "canvas")
	canvas.Set("width", 600)
	canvas.Set("height", 600)
	document.Get("body").Call("appendChild", canvas)

	observer := JSGameObserver{done, canvas}

	onKeyDown := func(this js.Value, p []js.Value) interface{} {
		keyCode := p[0].Get("keyCode").Int()
		if keyCode == 65 {
			toggleBox <- struct{}{}
		}
		if keyCode == 76 {
			toggleLetter <- struct{}{}
		}
		return js.Undefined()
	}
	js.Global().Call("addEventListener", "keydown", js.FuncOf(onKeyDown))

	go game.FlowLoop(ticker, toggleBox, toggleLetter, g, &observer)
	printScoreBoard(canvas, g)
	return js.Undefined()
}

func drawButtons(canvas js.Value, state string) {
	drawPlaceButton(canvas, state)
	drawTextButton(canvas, state)
}

func main() {
	fmt.Println("Go Web Assembly")
	js.Global().Set("run", js.FuncOf(run))
	select {}
}
