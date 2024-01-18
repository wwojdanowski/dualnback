package main

import (
	"fmt"
	"syscall/js"
)

func drawGrid(this js.Value, p []js.Value) interface{} {
	document := js.Global().Get("document")
	canvas := document.Call("createElement", "canvas")
	canvas.Set("width", 300)
	canvas.Set("height", 300)
	document.Get("body").Call("appendChild", canvas)

	context := canvas.Call("getContext", "2d")

	cellSize := 100
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			x := j * cellSize
			y := i * cellSize
			context.Call("fillRect", x, y, cellSize, cellSize)
		}
	}

	return js.Undefined()
}

func main() {
	fmt.Println("Go Web Assembly")
	js.Global().Set("drawGrid", js.FuncOf(drawGrid))
	select {}
}
