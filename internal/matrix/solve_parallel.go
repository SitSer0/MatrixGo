package matrix

import (
	"MatrixGo/internal/field"
	"MatrixGo/internal/vector"
	"errors"
	"runtime"
	"sync"
)

// SolveSystemParallel решает систему линейных уравнений параллельно
func SolveSystemParallel[T field.Field[T]](mat *Matrix[T], vec *vector.Vector[T]) (*vector.Vector[T], error) {
	if mat.Rows != mat.Cols {
		return nil, errors.New("матрица должна быть квадратной")
	}
	if mat.Rows != vec.Len() {
		return nil, errors.New("размер вектора не совпадает с размером матрицы")
	}

	n := mat.Rows
	M := mat.Clone()
	B := make([]T, n)
	copy(B, vec.Data)

	for i := 0; i < n; i++ {
		maxRow := i
		for j := i + 1; j < n; j++ {
			if !M.Data[j][i].Equal(M.Data[j][i].Zero()) {
				maxRow = j
				break
			}
		}

		if maxRow != i {
			M.Data[i], M.Data[maxRow] = M.Data[maxRow], M.Data[i]
			B[i], B[maxRow] = B[maxRow], B[i]
		}

		pivot := M.Data[i][i]
		if pivot.Equal(pivot.Zero()) {
			return nil, errors.New("матрица вырождена, решение невозможно")
		}

		for j := i; j < n; j++ {
			d, _ := M.Data[i][j].Div(pivot)
			M.Data[i][j] = d
		}
		d, _ := B[i].Div(pivot)
		B[i] = d

		var wg sync.WaitGroup
		numCPU := runtime.NumCPU()
		rowsPerGoroutine := (n - i - 1) / numCPU
		if rowsPerGoroutine < 1 {
			rowsPerGoroutine = 1
		}

		for startRow := i + 1; startRow < n; startRow += rowsPerGoroutine {
			endRow := startRow + rowsPerGoroutine
			if endRow > n {
				endRow = n
			}

			wg.Add(1)
			go func(start, end int) {
				defer wg.Done()
				for k := start; k < end; k++ {
					factor := M.Data[k][i]
					for j := i; j < n; j++ {
						M.Data[k][j] = M.Data[k][j].Sub(factor.Mul(M.Data[i][j]))
					}
					B[k] = B[k].Sub(factor.Mul(B[i]))
				}
			}(startRow, endRow)
		}
		wg.Wait()
	}

	x := make([]T, n)
	for i := n - 1; i >= 0; i-- {
		x[i] = B[i]
		for j := i + 1; j < n; j++ {
			x[i] = x[i].Sub(M.Data[i][j].Mul(x[j]))
		}
	}

	return vector.NewVector[T](x), nil
}
