package main

import (
	"errors"
	"fmt"
)

func dotProduct(x, y [][]float32) ([][]float32, error) {
	if len(x[0]) != len(y) {
		return nil, errors.New("Can't do matrix multiplication.")
	}

	out := make([][]float32, len(x))
	for i := 0; i < len(x); i += 1 {
		for j := 0; j < len(y); j += 1 {
			if len(out[i]) < 1 {
				out[i] = make([]float32, len(y))
			}
			out[i][j] += x[i][j] * y[j][i]
		}
	}
	return out, nil
}

func transpose(x [][]float32) [][]float32 {
	out := make([][]float32, len(x[0]))
	for i := 0; i < len(x); i += 1 {
		for j := 0; j < len(x[0]); j += 1 {
			out[j] = append(out[j], x[i][j])
		}
	}
	return out
}
func main() {
	X := [][]float32{
		[]float32{5, 8, 3},
		[]float32{9, 8, 5},
	}

	Y := [][]float32{
		[]float32{8, 4, 6},
		[]float32{2, 8, 1},
	}

	out, _ := dotProduct(X, transpose(Y))
	fmt.Println(out)
}
