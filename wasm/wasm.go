package main

import (
	"dualnback/game"
	"fmt"
	"syscall/js"
)

type JSGameObserver struct {
}

/*

const canvas = document.getElementById("myCanvas");
const ctx = canvas.getContext("2d");

ctx.font = "30px Arial";
ctx.fillText("Hello World", 10, 50);

*/

func drawGridWithItem(this js.Value, p []js.Value, newItem game.Item) interface{} {
	document := js.Global().Get("document")
	canvas := document.Call("createElement", "canvas")
	canvas.Set("width", 300)
	canvas.Set("height", 300)
	document.Get("body").Call("appendChild", canvas)

	context := canvas.Call("getContext", "2d")

	boxRow := newItem.Box / 3
	boxCol := newItem.Box % 3

	cellSize := 100
	cellMargin := 10
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			x := j * cellSize
			y := i * cellSize
			if i == boxRow && j == boxCol {
				context.Set("fillStyle", "green")
				context.Call("fillRect", x+cellMargin, y+cellMargin,
					cellSize-cellMargin, cellSize-cellMargin)
				context.Set("font", "30px Arial")
				context.Call("strokeText", "H", x+cellMargin, y+cellMargin)
				context.Set("fillStyle", "black")
			} else {
				context.Call("fillRect", x+cellMargin, y+cellMargin,
					cellSize-cellMargin, cellSize-cellMargin)
			}
		}
	}
	// letters := []rune{'A', 'B', 'C', 'D', 'E'}

	return js.Undefined()
}

func (o *JSGameObserver) NewSequence(g *game.Game, item game.Item) {
	// drawGridWithItem(o.s, 1, 1, 3, 3, boxStyle, itemStyle, newItem)
	// if g.IsReady() {
	// 	drawText(o.s, 10, 10, 50, 15, readyToBePressedStyle, "PLACE")
	// 	drawText(o.s, 18, 10, 50, 15, readyToBePressedStyle, "LETTER")
	// }
}

func (o *JSGameObserver) PauseForDecision(g *game.Game) {
	panic("not implemented") // TODO: Implement
}

func (o *JSGameObserver) EvalRound(g *game.Game) {
	panic("not implemented") // TODO: Implement
}

func (o *JSGameObserver) RoundFinished(g *game.Game) {
	panic("not implemented") // TODO: Implement
}

func (o *JSGameObserver) StateProcessed(g *game.Game) {
	panic("not implemented") // TODO: Implement
}

func (o *JSGameObserver) ToggleBox(g *game.Game) {
	panic("not implemented") // TODO: Implement
}

func (o *JSGameObserver) ToggleLetter(g *game.Game) {
	panic("not implemented") // TODO: Implement
}

func drawGrid(this js.Value, p []js.Value) interface{} {
	document := js.Global().Get("document")
	canvas := document.Call("createElement", "canvas")
	canvas.Set("width", 300)
	canvas.Set("height", 300)
	document.Get("body").Call("appendChild", canvas)

	context := canvas.Call("getContext", "2d")

	cellSize := 100
	cellMargin := 10
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

func drawTestGrid(this js.Value, p []js.Value) interface{} {
	newItem := game.Item{1, 1}
	return drawGridWithItem(this, p, newItem)
}

func main() {
	fmt.Println("Go Web Assembly")
	js.Global().Set("drawGrid", js.FuncOf(drawGrid))
	js.Global().Set("drawTestGrid", js.FuncOf(drawTestGrid))
	select {}
}
