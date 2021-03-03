package nrpa

import (
	"alda/entities"
	"alda/utils"
)

var (
	policyTmp [][]float64
)

// Data structure for each NRPA recursive call (level)
type Level struct {
	Policy      [][]float64
	BestRollout *Rollout
}

//Plays a new rollout
func (l *Level) PlayOut(t *entities.TSPTW) *Rollout {
	rollout := NewRollout(t)
	rollout.Do(l.Policy)
	l.BestRollout = rollout
	return rollout
}

// Adapt the level policy by increasing the probability of the current BestRollout
func (l *Level) AdaptPolicy() {
	var k int
	u := 0
	r := l.BestRollout
	utils.CopyPolicy(l.Policy, policyTmp) //copy level policy in a temporal policy copy

	for step := range r.legalMovesPerStep {
		v := r.Tour[step+1]
		moves := r.legalMovesPerStep[step]
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
