package nrpa

import (
	"math"
	"math/rand"

	"alda/entities"
	"alda/utils"
)

// Data structure for a NRPA Rollout
type Rollout struct {
	Tour              []int
	Length            int
	Score             float64
	Makespan          float64
	visited           []bool
	legalMovesPerStep [][]int
	moveProb          []float64
	cost              float64
	Violations        int
	t                 *entities.TSPTW
}

func NewRollout(t *entities.TSPTW) *Rollout {
	rollout := &Rollout{
		legalMovesPerStep: make([][]int, t.N-1),
		moveProb:          make([]float64, t.N),
		Tour:              make([]int, t.N+1),
		visited:           make([]bool, t.N),
		Length:            1,
		t:                 t,
	}
	for i := range rollout.legalMovesPerStep {
		rollout.legalMovesPerStep[i] = make([]int, 0, t.N)
	}
	return rollout
}

// Do a Rollout for a TSPTW with the given policy
func (r *Rollout) Do(policy [][]float64) float64 {
	u := 0 // starts at the depot

	for step := 0; r.Length < r.t.N; step++ {
		r.calculateLegalMoves(u, step)
		v := r.pickMove(u, step, policy)
		r.move(u, v)
		u = v
	}
	r.move(u, 0) // go back to the depot
	r.Score = -(float64(offset*r.Violations) + r.cost)

	return r.Score
}

// Reset the Rollout underling cases
func (r *Rollout) Reset() {
	r.Length = 1
	r.Score = -math.MaxFloat64
	r.Makespan = 0
	r.cost = 0
	r.Violations = 0
	for i := 1; i < r.t.N; i++ {
		r.visited[i] = false
	}
	for state := 0; state < r.t.N-1; state++ {
		r.legalMovesPerStep[state] = r.legalMovesPerStep[state][:0]
	}
}

// Calculates the legal moves that can be made from u node
func (r *Rollout) calculateLegalMoves(u, step int) {
	for v := 1; v < r.t.N; v++ {
		if !r.visited[v] {
			if r.Makespan+r.t.Distances[u][v] > r.t.WindowEnd[v] { // violation of the TW
				// take that move since there's no other path u->k->v faster than u->v
				r.legalMovesPerStep[step] = append(r.legalMovesPerStep[step], v) // that don't violate the TW constraint (Triangle Inequality)
			}
		}
	}
	if len(r.legalMovesPerStep[step]) == 0 {
		for v := 1; v < r.t.N; v++ {
			if !r.visited[v] {
				impossibleMove := false
				for k := 0; k < r.t.N && !impossibleMove; k++ {
					if !r.visited[v] {
						if r.Makespan <= r.t.WindowEnd[k] && r.Makespan+r.t.Distances[u][k] <= r.t.WindowEnd[k] && // is valid u->k
							(r.Makespan+r.t.Distances[u][v] > r.t.WindowEnd[k] || r.t.WindowStart[v] > r.t.WindowEnd[k]) { //take u->v make impossible to go to k
							impossibleMove = true
						}
					}
				}
				if !impossibleMove {
					r.legalMovesPerStep[step] = append(r.legalMovesPerStep[step], v)
				}
			}
		}
	}
	if len(r.legalMovesPerStep[step]) == 0 { // there's no moves that do not violate a TW constrain
		for v := 1; v < r.t.N; v++ {
			if !r.visited[v] {
				r.legalMovesPerStep[step] = append(r.legalMovesPerStep[step], v)
			}
		}
	}
}

// Picks a move from u to a random node according to the policy distribution
func (r *Rollout) pickMove(u, step int, policy [][]float64) int {
	z := float64(0)
	for i := range r.legalMovesPerStep[step] {
		v := r.legalMovesPerStep[step][i]
		r.moveProb[i] = utils.Exp(policy[u][v])
		z += r.moveProb[i]
	}
	idx := 0
	random := rand.Float64() * z
	probAcc := r.moveProb[idx]
	for probAcc < random {
		idx++
		probAcc += r.moveProb[idx]
	}
	return r.legalMovesPerStep[step][idx]
}

//Make a move from u to v
func (r *Rollout) move(u, v int) {
	r.Tour[r.Length] = v
	r.Length++
	r.visited[v] = true
	r.cost += r.t.Distances[u][v]
	r.Makespan = utils.Max(r.Makespan+r.t.Distances[u][v], r.t.WindowStart[v])
	if r.Makespan > r.t.WindowEnd[v] {
		r.Violations++
	}
}
