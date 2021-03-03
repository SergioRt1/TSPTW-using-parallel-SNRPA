package nrpa

import (
	"math"
	"math/rand"

	"alda/entities"
	"alda/utils"
)

// Data structure for a NRPA Rollout
type Rollout struct {
	Tour       []int
	Length     int
	Score      float64
	Makespan   float64
	cost       float64
	Violations int
	d          *PlayoutData
}

type PlayoutData struct {
	visited           []bool
	legalMovesPerStep [][]int
	moveProb          []float64
	t                 *entities.TSPTW
}

func (d *PlayoutData) Reset() {
	for i := 1; i < d.t.N; i++ {
		d.visited[i] = false
	}
	for state := 0; state < d.t.N-1; state++ {
		d.legalMovesPerStep[state] = d.legalMovesPerStep[state][:0]
	}

}

func NewPlayoutData(t *entities.TSPTW) *PlayoutData {
	d := &PlayoutData{
		legalMovesPerStep: make([][]int, t.N-1),
		moveProb:          make([]float64, t.N),
		visited:           make([]bool, t.N),
		t:                 t,
	}
	for i := range d.legalMovesPerStep {
		d.legalMovesPerStep[i] = make([]int, 0, t.N)
	}
	return d
}

func NewRollout(data *PlayoutData) *Rollout {
	return &Rollout{
		Tour:   make([]int, data.t.N+1),
		Length: 1,
		Score:  -math.MaxFloat64,
		d:      data,
	}
}

// Do a Rollout for a TSPTW with the given policy
func (r *Rollout) Do(policy [][]float64) float64 {
	u := 0 // starts at the depot

	for step := 0; r.Length < r.d.t.N; step++ {
		r.calculateLegalMoves(u, step)
		v := r.pickMove(u, step, policy)
		r.move(u, v)
		u = v
	}
	r.move(u, 0) // go back to the depot
	r.Score = -(float64(offset*r.Violations) + r.cost)

	return r.Score
}

// Calculates the legal moves that can be made from u node
func (r *Rollout) calculateLegalMoves(u, step int) {
	for v := 1; v < r.d.t.N; v++ {
		if !r.d.visited[v] {
			if r.Makespan+r.d.t.Distances[u][v] > r.d.t.WindowEnd[v] { // violation of the TW
				// take that move since there's no other path u->k->v faster than u->v
				r.d.legalMovesPerStep[step] = append(r.d.legalMovesPerStep[step], v) // that don't violate the TW constraint (Triangle Inequality)
			}
		}
	}
	if len(r.d.legalMovesPerStep[step]) == 0 {
		for v := 1; v < r.d.t.N; v++ {
			if !r.d.visited[v] {
				impossibleMove := false
				for k := 0; k < r.d.t.N && !impossibleMove; k++ {
					if !r.d.visited[v] {
						if r.Makespan <= r.d.t.WindowEnd[k] && r.Makespan+r.d.t.Distances[u][k] <= r.d.t.WindowEnd[k] && // is valid u->k
							(r.Makespan+r.d.t.Distances[u][v] > r.d.t.WindowEnd[k] || r.d.t.WindowStart[v] > r.d.t.WindowEnd[k]) { //take u->v make impossible to go to k
							impossibleMove = true
						}
					}
				}
				if !impossibleMove {
					r.d.legalMovesPerStep[step] = append(r.d.legalMovesPerStep[step], v)
				}
			}
		}
	}
	if len(r.d.legalMovesPerStep[step]) == 0 { // there's no moves that do not violate a TW constrain
		for v := 1; v < r.d.t.N; v++ {
			if !r.d.visited[v] {
				r.d.legalMovesPerStep[step] = append(r.d.legalMovesPerStep[step], v)
			}
		}
	}
}

// Picks a move from u to a random node according to the policy distribution
func (r *Rollout) pickMove(u, step int, policy [][]float64) int {
	z := float64(0)
	for i := range r.d.legalMovesPerStep[step] {
		v := r.d.legalMovesPerStep[step][i]
		r.d.moveProb[i] = utils.Exp(policy[u][v])
		z += r.d.moveProb[i]
	}
	idx := 0
	random := rand.Float64() * z
	probAcc := r.d.moveProb[idx]
	for probAcc < random {
		idx++
		probAcc += r.d.moveProb[idx]
	}
	return r.d.legalMovesPerStep[step][idx]
}

//Make a move from u to v
func (r *Rollout) move(u, v int) {
	r.Tour[r.Length] = v
	r.Length++
	r.d.visited[v] = true
	r.cost += r.d.t.Distances[u][v]
	r.Makespan = utils.Max(r.Makespan+r.d.t.Distances[u][v], r.d.t.WindowStart[v])
	if r.Makespan > r.d.t.WindowEnd[v] {
		r.Violations++
	}
}
