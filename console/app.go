package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	game "dualnback/game"

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

type SimpleGameObserver struct {
	s    tcell.Screen
	done chan bool
}

var defStyle = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
var boxStyle = tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack)
var itemStyle = tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorGreenYellow)
var correctStyle = tcell.StyleDefault.Background(tcell.ColorGreenYellow).Foreground(tcell.ColorBlack)
var wrongStyle = tcell.StyleDefault.Background(tcell.ColorRed).Foreground(tcell.ColorWhite)

var notReadyToBePressedStyle = tcell.StyleDefault.Background(tcell.ColorGray).Foreground(tcell.ColorWhite)
var readyToBePressedStyle = tcell.StyleDefault.Background(tcell.ColorDarkGray).Foreground(tcell.ColorWhite)
var toggledStyle = tcell.StyleDefault.Background(tcell.ColorBlue).Foreground(tcell.ColorWhite)

func (o *SimpleGameObserver) NewSequence(g *game.Game, newItem game.Item) {
	drawGridWithItem(o.s, 1, 1, 3, 3, boxStyle, itemStyle, newItem)
	if g.IsReady() {
		drawText(o.s, 10, 10, 50, 15, readyToBePressedStyle, "PLACE")
		drawText(o.s, 18, 10, 50, 15, readyToBePressedStyle, "LETTER")
	}
}

func (o *SimpleGameObserver) PauseForDecision(g *game.Game) {
	drawGrid(o.s, 1, 1, 3, 3, boxStyle)
}

func (o *SimpleGameObserver) EvalRound(g *game.Game) {

	if g.IsReady() {
		if g.LastResult.Box {
			drawText(o.s, 10, 10, 50, 15, correctStyle, "PLACE")
		} else {
			drawText(o.s, 10, 10, 50, 15, wrongStyle, "PLACE")
		}

		if g.LastResult.Letter {
			drawText(o.s, 18, 10, 50, 15, correctStyle, "LETTER")
		} else {
			drawText(o.s, 18, 10, 50, 15, wrongStyle, "LETTER")
		}
	}

	printScoreBoard(o.s, 10, 1, 50, 15, g.N, g.Score, g.Round, g.MaxRounds, defStyle)
}

func (o *SimpleGameObserver) RoundFinished(g *game.Game) {
	drawText(o.s, 10, 10, 50, 15, readyToBePressedStyle, "PLACE")
	drawText(o.s, 18, 10, 50, 15, readyToBePressedStyle, "LETTER")
}

func (o *SimpleGameObserver) ToggleBox(g *game.Game) {
	if g.IsBoxToggled() {
		drawText(o.s, 10, 10, 50, 15, toggledStyle, "PLACE")
	} else {
		drawText(o.s, 10, 10, 50, 15, readyToBePressedStyle, "PLACE")
	}
	o.s.Sync()
}

func (o *SimpleGameObserver) ToggleLetter(g *game.Game) {
	if g.IsLetterToggled() {
		drawText(o.s, 18, 10, 50, 15, toggledStyle, "LETTER")
	} else {
		drawText(o.s, 18, 10, 50, 15, readyToBePressedStyle, "LETTER")
	}
	o.s.Sync()
}

func (o *SimpleGameObserver) StateProcessed(g *game.Game) {
	o.s.Sync()
	if g.IsDone() {
		printScoreBoard(o.s, 10, 1, 50, 15, g.N, g.Score, g.Round, g.MaxRounds, defStyle)
		o.s.Sync()
		o.done <- true
		return
	}
}

func controlLoop(s tcell.Screen, done chan bool, ticker *time.Ticker, toggleBox chan struct{}, toggleLetter chan struct{}) {
	for {
		s.Show()
		ev := s.PollEvent()
		select {
		case <-done:
			ticker.Stop()
			close(toggleBox)
			close(toggleLetter)
			drawBox(s, 0, 0, 60, 12, defStyle, "We're done!")
			s.Sync()
		default:
			switch ev := ev.(type) {
			case *tcell.EventResize:
				s.Sync()
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
					return
				} else if ev.Key() == tcell.KeyCtrlL {
					s.Sync()
				} else if ev.Rune() == 'A' || ev.Rune() == 'a' {
					toggleBox <- struct{}{}
				} else if ev.Rune() == 'L' || ev.Rune() == 'l' {
					toggleLetter <- struct{}{}
				}
			case *tcell.EventMouse:

				switch ev.Buttons() {
				case tcell.Button1, tcell.Button2:

				case tcell.ButtonNone:
				}
			}
		}
	}
}

func main() {

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

	nbackPtr := flag.Int("n", 2, "N-back")
	roundsPtr := flag.Int("rounds", 2, "Max rounds")

	flag.Parse()
	g := game.NewGame(*nbackPtr, *roundsPtr)

	drawBox(s, 0, 0, 60, 12, defStyle, "")
	drawGrid(s, 1, 1, 3, 3, boxStyle)
	printScoreBoard(s, 10, 1, 50, 15, g.N, g.Score, g.Round, g.MaxRounds, defStyle)

	drawText(s, 10, 10, 50, 15, notReadyToBePressedStyle, "PLACE")
	drawText(s, 18, 10, 50, 15, notReadyToBePressedStyle, "LETTER")

	ticker := time.NewTicker(500 * time.Millisecond)
	toggleBox := make(chan struct{})
	toggleLetter := make(chan struct{})
	done := make(chan bool)

	observer := SimpleGameObserver{s, done}

	go game.FlowLoop(ticker, toggleBox, toggleLetter, g, &observer)

	quit := func() {
		maybePanic := recover()
		s.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}
	defer quit()
	controlLoop(s, done, ticker, toggleBox, toggleLetter)
}
