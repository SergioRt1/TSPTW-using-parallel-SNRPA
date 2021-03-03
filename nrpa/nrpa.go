package nrpa

import (
	"alda/entities"
	"alda/utils"
	"math"
)

const (
	offset = 1000000
	alpha  = 1.0
)

type nrpa struct {
	NInter       int
	DataPerLevel []*Level
	t            *entities.TSPTW
}

func NewNRPA(tsptw *entities.TSPTW, levels, nIter int) *nrpa {
	return &nrpa{
		NInter:       nIter,
		DataPerLevel: make([]*Level, levels+1),
		t:            tsptw,
	}
}

func (n *nrpa) StableNRPA(level int, nLevel *Level, policy [][]float64, factor int) *Rollout {
	if level == 0 {
		return nLevel.PlayOut(n.t)
	}
	nLevel.BestRollout.Score = -math.MaxFloat64
	nextLevel := n.DataPerLevel[level-1]
	utils.CopyPolicy(policy, nLevel.Policy)
	if level == 1 {
		for i := 0; i < factor; i++ {
			_ = n.StableNRPA(level-1, nextLevel, nLevel.Policy, factor)
			if nextLevel.BestRollout.Score >= nLevel.BestRollout.Score {
				nLevel.BestRollout, nextLevel.BestRollout = nextLevel.BestRollout, nLevel.BestRollout
			}
		}
	} else {
		for i := 0; i < n.NInter; i++ {
			_ = n.StableNRPA(level-1, nextLevel, nLevel.Policy, factor)
			if nextLevel.BestRollout.Score >= nLevel.BestRollout.Score {
				nLevel.BestRollout, nextLevel.BestRollout = nextLevel.BestRollout, nLevel.BestRollout
			}
			nLevel.AdaptPolicy()
		}
	}

	return nLevel.BestRollout
}

func (n *nrpa) PreAllocate() [][]float64 {
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
