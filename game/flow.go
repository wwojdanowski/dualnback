package game

import (
	"time"
)

const (
	NewSequenceState = iota
	PauseForDecisionState
	EvalRoundState

	NewSequenceStateGameReady
	PauseForDecisionStateGameReady
	EvalRoundStateGameReady
	PostRoundStateGameReady
)

func FlowLoop(ticker *time.Ticker, toggleBox <-chan struct{}, toggleLetter <-chan struct{}, g *Game, observer GameObserver) {
	state := NewSequenceState
	for {
		select {
		case <-ticker.C:
			switch state {
			case NewSequenceState:
				newItem := MakeRandomItem()
				g.NextSequence(newItem)
				observer.NewSequence(g, newItem)
				if g.IsReady() {
					state = PauseForDecisionStateGameReady
				} else {
					state = PauseForDecisionState
				}
				ticker.Reset(2000 * time.Millisecond)
			case PauseForDecisionState:
				observer.PauseForDecision(g)
				state = EvalRoundState
			case EvalRoundState:
				observer.EvalRound(g)
				state = NewSequenceState
				ticker.Reset(1000 * time.Millisecond)
			case NewSequenceStateGameReady:
				newItem := MakeRandomItem()
				g.NextSequence(newItem)
				observer.NewSequence(g, newItem)
				state = PauseForDecisionStateGameReady
				ticker.Reset(2000 * time.Millisecond)
			case PauseForDecisionStateGameReady:
				observer.PauseForDecision(g)
				state = EvalRoundStateGameReady
			case EvalRoundStateGameReady:
				g.EvalRound()
				observer.EvalRound(g)
				state = PostRoundStateGameReady
				ticker.Reset(1000 * time.Millisecond)
			case PostRoundStateGameReady:
				observer.RoundFinished(g)
				state = NewSequenceStateGameReady
			}
			observer.StateProcessed(g)

		case <-toggleBox:
			if state == PauseForDecisionStateGameReady || state == EvalRoundStateGameReady {
				g.ToggleBox()
				observer.ToggleBox(g)

			}
		case <-toggleLetter:
			if state == PauseForDecisionStateGameReady || state == EvalRoundStateGameReady {
				g.ToggleLetter()
				observer.ToggleLetter(g)
			}
		}
	}
}
