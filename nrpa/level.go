package nrpa

import (
	"alda/utils"
)

var (
	policyTmp [][]float64
)

// Data structure for each NRPA recursive call (level)
type Level struct {
	Policy            [][]float64
	BestRollout       *Rollout
	LegalMovesPerStep [][]int
}

// Adapt the level policy by increasing the probability of the current BestRollout
func (l *Level) AdaptPolicy() {
	var k int
	u := 0
	r := l.BestRollout
	utils.CopyPolicy(l.Policy, policyTmp) //copy level policy in a temporal policy copy

	for step := range l.LegalMovesPerStep {
		v := r.Tour[step+1]
		moves := l.LegalMovesPerStep[step]
		l.Policy[u][v] += alpha
		z := 0.0
		for m := range moves {
			k = moves[m]
			z += utils.Exp(policyTmp[u][k])
		}
		for m := range moves {
			k = moves[m]
			l.Policy[u][k] -= alpha * utils.Exp(policyTmp[u][k]) / z
		}
		u = v
	}
}
