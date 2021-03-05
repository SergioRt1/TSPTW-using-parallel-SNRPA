package nrpa

import (
	"context"
	"math"
	"sync"

	"alda/cli"
	"alda/entities"
	"alda/utils"
)

const (
	offset = 1000000
	alpha  = 1.0
)

type StaticData struct {
	bestMoves [][]int
	policyTmp [][]float64
}

type NRPA struct {
	NInter              int
	DataPerLevel        []*Level
	StabilizationFactor int
	Levels              int
	t                   *entities.TSPTW
	StaticData          *StaticData
	Actors              []*playoutActor
}

func NewNRPA(tsptw *entities.TSPTW, config *cli.Config, data *StaticData) *NRPA {
	return &NRPA{
		NInter:              config.NIter,
		DataPerLevel:        make([]*Level, config.Levels),
		Actors:              make([]*playoutActor, config.PActors),
		t:                   tsptw,
		Levels:              config.Levels,
		StabilizationFactor: config.StabilizationFactor,
		StaticData:          data,
	}
}

//Run Parallel Stable NRPA with actors to compute the leaves
func (n *NRPA) RunConcurrent(ctx context.Context, levels int, t *entities.TSPTW, out chan *Rollout, wg *sync.WaitGroup) {
	defer wg.Done()
	var bestRollout *Rollout

	done := make(chan *Rollout)
	for i := range n.Actors {
		n.Actors[i] = StartPlayoutActor(ctx, t)
	}
	policy := n.PreAllocate()
	go func() {
		done <- n.StableNRPA(levels-1, n.DataPerLevel[levels-1], policy)
	}()
	select {
	case <-ctx.Done():
		bestRollout = n.FindCurrentBest()
	case bestRollout = <-done:
	}
	out <- bestRollout
}

func (n *NRPA) StableNRPA(level int, nLevel *Level, policy [][]float64) *Rollout {
	nLevel.BestRollout.Score = -math.MaxFloat64
	if level == 1 {
		n.concurrentPlayout(policy, nLevel)
	} else {
		utils.CopyPolicy(policy, nLevel.Policy)
		nextLevel := n.DataPerLevel[level-1]
		for i := 0; i < n.NInter; i++ {
			_ = n.StableNRPA(level-1, nextLevel, nLevel.Policy)
			if nextLevel.BestRollout.Score >= nLevel.BestRollout.Score {
				nLevel.BestRollout, nextLevel.BestRollout = nextLevel.BestRollout, nLevel.BestRollout
				utils.CopyMoves(nextLevel.LegalMovesPerStep, nLevel.LegalMovesPerStep)
			}
			nLevel.AdaptPolicy(n.StaticData.policyTmp)
		}
	}

	return nLevel.BestRollout
}

func (n *NRPA) concurrentPlayout(policy [][]float64, nLevel *Level) {
	chOut := make(chan *Message, n.StabilizationFactor)
	var wg sync.WaitGroup

	wg.Add(n.StabilizationFactor)
	go func() {
		wg.Wait()
		close(chOut)
	}()
	for i := 0; i < n.StabilizationFactor; i++ {
		n.Actors[i%len(n.Actors)].Playout(policy, chOut, &wg)
	}

	for message := range chOut {
		if message.Rollout.Score >= nLevel.BestRollout.Score {
			nLevel.BestRollout = message.Rollout
			n.StaticData.bestMoves = message.LegalMovesPerStep
		}
	}
	utils.CopyMoves(n.StaticData.bestMoves, nLevel.LegalMovesPerStep)
}

func (n *NRPA) PreAllocate() [][]float64 {
	for i := range n.DataPerLevel {
		level := &Level{
			Policy:            make([][]float64, n.t.N),
			BestRollout:       &Rollout{},
			LegalMovesPerStep: make([][]int, n.t.N-1),
		}
		for j := range level.Policy {
			level.Policy[j] = make([]float64, n.t.N)
		}
		for j := range level.LegalMovesPerStep {
			level.LegalMovesPerStep[j] = make([]int, n.t.N)
		}
		n.DataPerLevel[i] = level
	}
	policy := make([][]float64, n.t.N) // Policy used to pass accumulated knowledge to the lower levels
	for i := range policy {
		policy[i] = make([]float64, n.t.N)
	}
	return policy
}

func (n *NRPA) FindCurrentBest() *Rollout {
	bestRollout := &Rollout{Score: -math.MaxFloat64}
	for i := n.Levels - 1; i >= 0; i-- {
		levelBest := *n.DataPerLevel[i].BestRollout
		if levelBest.Score > bestRollout.Score && levelBest.Length == n.t.N+1 {
			bestRollout = &levelBest
		}
	}
	return bestRollout
}

func NewStaticData(t *entities.TSPTW) *StaticData {
	d := &StaticData{
		bestMoves: make([][]int, t.N),     // Legal best moves instance used as temporary variable for copying
		policyTmp: make([][]float64, t.N), // Policy instance used as temporary variable for copying
	}
	for i := range d.policyTmp {
		d.policyTmp[i] = make([]float64, t.N)
	}
	for i := range d.bestMoves {
		d.bestMoves[i] = make([]int, t.N)
	}
	return d
}
