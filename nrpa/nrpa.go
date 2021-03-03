package nrpa

import (
	"context"
	"math"

	"alda/entities"
	"alda/utils"
)

const (
	offset = 1000000
	alpha  = 1.0
)

type NRPA struct {
	NInter              int
	DataPerLevel        []*Level
	StabilizationFactor int
	Levels              int
	t                   *entities.TSPTW
}

func NewNRPA(tsptw *entities.TSPTW, levels, nIter, factor int) *NRPA {
	return &NRPA{
		NInter:              nIter,
		DataPerLevel:        make([]*Level, levels+1),
		t:                   tsptw,
		Levels:              levels,
		StabilizationFactor: factor,
	}
}

func (n *NRPA) StableNRPA(ctx context.Context, level int, nLevel *Level, policy [][]float64) *Rollout {
	if level == 0 {
		return nLevel.PlayOut(n.t)
	}
	nLevel.BestRollout.Score = -math.MaxFloat64
	nextLevel := n.DataPerLevel[level-1]
	utils.CopyPolicy(policy, nLevel.Policy)
	if level == 1 {
		for i := 0; i < n.StabilizationFactor; i++ {
			_ = n.StableNRPA(ctx, level-1, nextLevel, nLevel.Policy)
			if nextLevel.BestRollout.Score >= nLevel.BestRollout.Score {
				nLevel.BestRollout, nextLevel.BestRollout = nextLevel.BestRollout, nLevel.BestRollout
			}
		}
	} else {
		for i := 0; i < n.NInter; i++ {
			_ = n.StableNRPA(ctx, level-1, nextLevel, nLevel.Policy)
			if nextLevel.BestRollout.Score >= nLevel.BestRollout.Score {
				nLevel.BestRollout, nextLevel.BestRollout = nextLevel.BestRollout, nLevel.BestRollout
			}
			nLevel.AdaptPolicy()
		}
	}

	return nLevel.BestRollout
}

func (n *NRPA) PreAllocate() [][]float64 {
	for i := range n.DataPerLevel {
		level := &Level{
			Policy:      make([][]float64, n.t.N),
			BestRollout: &Rollout{},
		}
		for j := range level.Policy {
			level.Policy[j] = make([]float64, n.t.N)
		}
		n.DataPerLevel[i] = level
	}
	policy := make([][]float64, n.t.N) // Policy used to pass accumulated knowledge to the lower levels
	for i := range policy {
		policy[i] = make([]float64, n.t.N)
	}
	policyTmp = make([][]float64, n.t.N) // Policy instance used as temporary variable for copying
	for i := range policyTmp {
		policyTmp[i] = make([]float64, n.t.N)
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
