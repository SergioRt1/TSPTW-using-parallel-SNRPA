package nrpa

import (
	"alda/cli"
	"alda/entities"
	"context"
	"sync"
)

type NrpaActor struct {
	data *StaticData
	chIn chan action
	ctx  context.Context
	t    *entities.TSPTW
}

func StartNRPAActor(ctx context.Context, t *entities.TSPTW) *NrpaActor {
	a := &NrpaActor{
		data: NewStaticData(t),
		chIn: make(chan action),
		ctx:  ctx,
		t:    t,
	}
	go a.Loop(ctx)
	return a
}

func (a *NrpaActor) RunNRPA(config *cli.Config, chOut chan *Rollout, wg *sync.WaitGroup) {
	a.chIn <- func() {
		nrpaInstance := NewNRPA(a.t, config, a.data)
		nrpaInstance.RunConcurrent(a.ctx, config.Levels, a.t, chOut, wg)
	}
}

func (a *NrpaActor) Loop(ctx context.Context) {
	for {
		select {
		case do := <-a.chIn:
			do()
		case <-ctx.Done():
			return
		}
	}
}
