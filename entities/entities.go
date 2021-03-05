package entities

type TSPTW struct {
	N           int
	WindowStart []float64
	WindowEnd   []float64
	Distances   [][]float64
	Delta       float64
	Min         float64
}

func (t *TSPTW) Bias(i, j int) float64 {
	return 10 * (t.Distances[i][j] - t.Min) / t.Delta
}
