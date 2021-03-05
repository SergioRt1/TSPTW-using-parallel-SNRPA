package nrpa

import (
	"alda/entities"
	"alda/utils"
)

// Data structure for each NRPA recursive call (level)
type Level struct {
	Policy            [][]float64
	BestRollout       *Rollout
	LegalMovesPerStep [][]int
	moveProb          []float64
	t                 *entities.TSPTW
}

// Adapt the level policy by increasing the probability of the current BestRollout
func (l *Level) AdaptPolicy(policyTmp [][]float64) {
	var k int
	u := 0
	r := l.BestRollout
	utils.CopyPolicy(l.Policy, policyTmp) //copy level policy in a temporal policy copy

	for step := range l.LegalMovesPerStep {
		v := r.Tour[step+1]
		moves := l.LegalMovesPerStep[step]
		z := 0.0
		for m := range moves {
			k = moves[m]
			l.moveProb[m] = utils.Exp(policyTmp[u][k]/tau) + l.t.Bias(u, k)
			z += l.moveProb[m]
		}
		var bm float64
		for m := range moves {
			k = moves[m]
			if k == v {
				bm = 1
			} else {
				bm = 0
			}
			l.Policy[u][k] -= alpha / tau * (l.moveProb[m]/z - bm)
		}
		u = v
	}
}
