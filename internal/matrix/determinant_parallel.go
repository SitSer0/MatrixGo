package matrix

import (
	"runtime"
	"sync"
)

func (m *Matrix[T]) DeterminantParallel() T {
	if m.Rows != m.Cols {
		return m.Data[0][0].Zero()
	}

	mat := m.Clone()
	n := mat.Rows
	det := mat.Data[0][0].One()

	for i := 0; i < n; i++ {
		if mat.Data[i][i].Equal(mat.Data[i][i].Zero()) {
			swapped := false
			for j := i + 1; j < n; j++ {
				if !mat.Data[j][i].Equal(mat.Data[j][i].Zero()) {
					mat.Data[i], mat.Data[j] = mat.Data[j], mat.Data[i]
					det = det.Neg()
					swapped = true
					break
				}
			}
			if !swapped {
				return mat.Data[0][0].Zero()
			}
		}

		det = det.Mul(mat.Data[i][i])

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
				for j := start; j < end; j++ {
					if !mat.Data[j][i].Equal(mat.Data[j][i].Zero()) {
						factor, _ := mat.Data[j][i].Div(mat.Data[i][i])
						for k := i; k < n; k++ {
							mat.Data[j][k] = mat.Data[j][k].Sub(factor.Mul(mat.Data[i][k]))
						}
					}
				}
			}(startRow, endRow)
		}
		wg.Wait()
	}

	return det
}
