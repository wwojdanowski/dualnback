package main

import (
	"fmt"
	"log"
	"time"

	game "dualnback/wasm"

	tcell "github.com/gdamore/tcell/v2"
)

func drawText(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	row := y1
	col := x1
	for _, r := range []rune(text) {
		s.SetContent(col, row, r, nil, style)
		col++
		if col >= x2 {
			row++
			col = x1
		}
		if row > y2 {
			break
		}
	}
}

func drawBox(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	if y2 < y1 {
		y1, y2 = y2, y1
	}
	if x2 < x1 {
		x1, x2 = x2, x1
	}

	// Fill background
	for row := y1; row <= y2; row++ {
		for col := x1; col <= x2; col++ {
			s.SetContent(col, row, ' ', nil, style)
		}
	}

	// Draw borders
	for col := x1; col <= x2; col++ {
		s.SetContent(col, y1, tcell.RuneHLine, nil, style)
		s.SetContent(col, y2, tcell.RuneHLine, nil, style)
	}
	for row := y1 + 1; row < y2; row++ {
		s.SetContent(x1, row, tcell.RuneVLine, nil, style)
		s.SetContent(x2, row, tcell.RuneVLine, nil, style)
	}

	// Only draw corners if necessary
	if y1 != y2 && x1 != x2 {
		s.SetContent(x1, y1, tcell.RuneULCorner, nil, style)
		s.SetContent(x2, y1, tcell.RuneURCorner, nil, style)
		s.SetContent(x1, y2, tcell.RuneLLCorner, nil, style)
		s.SetContent(x2, y2, tcell.RuneLRCorner, nil, style)
	}

	drawText(s, x1+1, y1+1, x2-1, y2-1, style, text)
}

func drawGrid(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style) {
	if y2 < y1 {
		y1, y2 = y2, y1
	}
	if x2 < x1 {
		x1, x2 = x2, x1
	}

	for row := y1; row <= y2; row++ {
		for col := x1; col <= x2; col++ {
			s.SetContent(col*2, row*2, tcell.RuneBlock, nil, style)
		}
	}
}

func drawGridWithItem(s tcell.Screen, x1, y1, x2, y2 int,
	style tcell.Style,
	itemStyle tcell.Style,
	item game.Item) {
	if y2 < y1 {
		y1, y2 = y2, y1
	}
	if x2 < x1 {
		x1, x2 = x2, x1
	}

	boxRow := item.Box/3 + 1
	boxCol := item.Box%3 + 1

	for row := y1; row <= y2; row++ {

		for col := x1; col <= x2; col++ {
			s.SetContent(col*2, row*2, tcell.RuneBlock, nil, style)
		}
	}

	letters := []rune{'A', 'B', 'C', 'D', 'E'}

	s.SetContent(boxCol*2, boxRow*2, letters[item.Letter], nil, itemStyle)

}

func printScoreBoard(s tcell.Screen, x1, y1, x2, y2 int, n, score, rounds, maxRounds int, style tcell.Style) {
	drawText(s, x1, y1, x2, y2, style, fmt.Sprintf("N: %d | score: %d | rounds: %d/%d", n, score, rounds, maxRounds))
}

func drawInputStatus(s tcell.Screen, x1, y1, x2, y2 int, result game.Result, style tcell.Style) {
}

func main() {
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	boxStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack)
	itemStyle := tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorGreenYellow)
	correctStyle := tcell.StyleDefault.Background(tcell.ColorGreenYellow).Foreground(tcell.ColorBlack)
	wrongStyle := tcell.StyleDefault.Background(tcell.ColorRed).Foreground(tcell.ColorWhite)

	notReadyToBePressedStyle := tcell.StyleDefault.Background(tcell.ColorGray).Foreground(tcell.ColorWhite)
	readyToBePressedStyle := tcell.StyleDefault.Background(tcell.ColorDarkGray).Foreground(tcell.ColorWhite)
	// toggledStyle := tcell.StyleDefault.Background(tcell.ColorBlue).Foreground(tcell.ColorWhite)

	// Initialize screen
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	s.SetStyle(defStyle)
	s.EnableMouse()
	s.EnablePaste()
	s.Clear()

	g := game.NewGame(2, 50)

	drawBox(s, 0, 0, 60, 12, defStyle, "")
	drawGrid(s, 1, 1, 3, 3, boxStyle)
	printScoreBoard(s, 10, 1, 50, 15, g.N, g.Score, g.Round, g.MaxRounds, defStyle)

	drawText(s, 10, 10, 50, 15, notReadyToBePressedStyle, "PLACE")
	drawText(s, 18, 10, 50, 15, notReadyToBePressedStyle, "LETTER")

	ticker := time.NewTicker(3000 * time.Millisecond)

	go func() {
		eval := false
		for {
			select {
			case <-ticker.C:
				// printScoreBoard(s, 10, 1, 50, 15, g.N, g.Score, g.Round, g.MaxRounds, defStyle)
				if eval {
					if g.IsReady() {
						g.EvalRound()
						if g.LastResult.Box {
							drawText(s, 10, 10, 50, 15, correctStyle, "PLACE")
						} else {
							drawText(s, 10, 10, 50, 15, wrongStyle, "PLACE")
						}

						if g.LastResult.Letter {
							drawText(s, 18, 10, 50, 15, correctStyle, "LETTER")
						} else {
							drawText(s, 18, 10, 50, 15, wrongStyle, "LETTER")
						}
					}
					eval = false
					drawGrid(s, 1, 1, 3, 3, boxStyle)
				} else {
					newItem := game.MakeRandomItem()
					g.NextSequence(newItem)
					drawGridWithItem(s, 1, 1, 3, 3, boxStyle, itemStyle, newItem)
					if g.IsReady() {
						drawText(s, 10, 10, 50, 15, readyToBePressedStyle, "PLACE")
						drawText(s, 18, 10, 50, 15, readyToBePressedStyle, "LETTER")
					}
					eval = true
				}
				s.Sync()
				if g.IsDone() {
					printScoreBoard(s, 10, 1, 50, 15, g.N, g.Score, g.Round, g.MaxRounds, defStyle)
					s.Sync()
					return
				}
			}
		}
	}()

	quit := func() {
		// You have to catch panics in a defer, clean up, and
		// re-raise them - otherwise your application can
		// die without leaving any diagnostic trace.
		maybePanic := recover()
		s.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}
	defer quit()

	// Here's how to get the screen size when you need it.
	// xmax, ymax := s.Size()

	// Here's an example of how to inject a keystroke where it will
	// be picked up by the next PollEvent call.  Note that the
	// queue is LIFO, it has a limited length, and PostEvent() can
	// return an error.
	// s.PostEvent(tcell.NewEventKey(tcell.KeyRune, rune('a'), 0))

	// Event loop
	for {
		// Update screen
		s.Show()

		// Poll event
		ev := s.PollEvent()

		// Process event
		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				return
			} else if ev.Key() == tcell.KeyCtrlL {
				s.Sync()
			} else if ev.Rune() == 'C' || ev.Rune() == 'c' {
				s.Clear()
			}
		case *tcell.EventMouse:

			switch ev.Buttons() {
			case tcell.Button1, tcell.Button2:

			case tcell.ButtonNone:
			}
		}
	}
}
