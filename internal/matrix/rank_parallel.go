package matrix

import (
	"runtime"
	"sync"
)

func (m *Matrix[T]) RankParallel() int {
	mat := m.Clone()
	rank := 0
	rowCount := mat.Rows
	colCount := mat.Cols
	h := 0
	k := 0

	for h < rowCount && k < colCount {
		i_max := h
		for i := h + 1; i < rowCount; i++ {
			if !mat.Data[i][k].Equal(mat.Data[i][k].Zero()) {
				i_max = i
				break
			}
		}

		if mat.Data[i_max][k].Equal(mat.Data[i_max][k].Zero()) {
			k++
			continue
		}

		if i_max != h {
			mat.Data[h], mat.Data[i_max] = mat.Data[i_max], mat.Data[h]
		}

		var wg sync.WaitGroup
		numCPU := runtime.NumCPU()
		rowsPerGoroutine := (rowCount - h - 1) / numCPU
		if rowsPerGoroutine < 1 {
			rowsPerGoroutine = 1
		}

		for startRow := h + 1; startRow < rowCount; startRow += rowsPerGoroutine {
			endRow := startRow + rowsPerGoroutine
			if endRow > rowCount {
				endRow = rowCount
			}

			wg.Add(1)
			go func(start, end int) {
				defer wg.Done()
				for i := start; i < end; i++ {
					if !mat.Data[i][k].Equal(mat.Data[i][k].Zero()) {
						factor, _ := mat.Data[i][k].Div(mat.Data[h][k])
						for j := k; j < colCount; j++ {
							mat.Data[i][j] = mat.Data[i][j].Sub(factor.Mul(mat.Data[h][j]))
						}
					}
				}
			}(startRow, endRow)
		}
		wg.Wait()

		rank++
		h++
		k++
	}

	return rank
}
