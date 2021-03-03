package utils

// Exp returns an approximation of e**x, the base-e exponential of x.
func Exp(x float64) float64 {
	x = 1.0 + x/1024
	x *= x
	x *= x
	x *= x
	x *= x
	x *= x
	x *= x
	x *= x
	x *= x
	x *= x
	x *= x

	return x
}

func Max(x, y float64) float64 {
	if x > y {
		return x
	}
	return y
}

// Copy the content of a inside b
func CopyPolicy(a, b [][]float64) {
	for i := range a {
		for j := range a[i] {
			b[i][j] = a[i][j]
		}
	}
}

// Copy the content of a inside b
func CopyMoves(a, b [][]int) {
	for i := range a {
		b[i] = b[i][:0]
		for j := range a[i] {
			b[i] = append(b[i], a[i][j])
		}
	}
}
