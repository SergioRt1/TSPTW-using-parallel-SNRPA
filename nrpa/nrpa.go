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
	Actors              []*actor
}

func NewNRPA(tsptw *entities.TSPTW, levels, nIter, factor int) *NRPA {
	return &NRPA{
		NInter:              nIter,
		DataPerLevel:        make([]*Level, levels),
		Actors:              make([]*actor, factor),
		t:                   tsptw,
		Levels:              levels,
		StabilizationFactor: factor,
		StaticData: &StaticData{
			bestMoves: make([][]int, tsptw.N),
			policyTmp: make([][]float64, tsptw.N),
		},
	}
}

//Run Parallel Stable NRPA with actors
func (n *NRPA) RunConcurrent(ctx context.Context, config *cli.Config, t *entities.TSPTW, out chan *Rollout, wg *sync.WaitGroup) {
	defer wg.Done()
	done := make(chan *Rollout)
	for i := 0; i < config.StabilizationFactor; i++ {
		n.Actors[i] = StartActor(ctx, t)
	}
	policy := n.PreAllocate()
	var bestRollout *Rollout
	go func() {
		done <- n.StableNRPA(config.Levels-1, n.DataPerLevel[config.Levels-1], policy)
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
		n.Actors[i].Playout(policy, chOut, &wg)
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
	// Policy instance used as temporary variable for copying
	for i := range n.StaticData.policyTmp {
		n.StaticData.policyTmp[i] = make([]float64, n.t.N)
	}
	// Legal best moves instance used as temporary variable for copying
	for i := range n.StaticData.bestMoves {
		n.StaticData.bestMoves[i] = make([]int, n.t.N)
	}
	return policy
}

func (n *NRPA) FindCurrentBest() *Rollout {
	bestRollout := &Rollout{Score: -math.MaxFloat64}
	for i := n.Levels; i >= 0; i-- {
		levelBest := *n.DataPerLevel[i].BestRollout
		if levelBest.Score > bestRollout.Score && levelBest.Length == n.t.N+1 {
			bestRollout = &levelBest
		}
	}
	return bestRollout
}
