package nrpa

import (
	"alda/entities"
	"context"
	"sync"
)

type action func()

type actor struct {
	t    *entities.TSPTW
	data *PlayoutData
	chIn chan action
}
type Message struct {
	Rollout           *Rollout
	LegalMovesPerStep [][]int
}

func StartActor(ctx context.Context, t *entities.TSPTW) *actor {
	a := &actor{
		t:    t,
		chIn: make(chan action),
		data: NewPlayoutData(t),
	}
	go a.Loop(ctx)
	return a
}

func (a *actor) Playout(policy [][]float64, chOut chan *Message, wg *sync.WaitGroup) {
	a.chIn <- func() {
		defer wg.Done()
		a.data.Reset()
		r := NewRollout(a.data)
		r.Do(policy)
		chOut <- &Message{
			Rollout:           r,
			LegalMovesPerStep: a.data.legalMovesPerStep,
		}
	}
}

func (a *actor) Loop(ctx context.Context) {
	for {
		select {
		case do := <-a.chIn:
			do()
		case <-ctx.Done():
			return
		}
	}
}
