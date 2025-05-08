package matrix

import (
	"errors"
	"runtime"
	"sync"
)

// InverseParallel вычисляет обратную матрицу через параллельное гауссово исключение
func (m *Matrix[T]) InverseParallel() (*Matrix[T], error) {
	if m.Rows != m.Cols {
		return nil, errors.New("матрица должна быть квадратной")
	}
	n := m.Rows
	A := m.Clone()
	I := IdentityMatrix[T](n, m.Data[0][0].Zero(), m.Data[0][0].One())

	for i := 0; i < n; i++ {
		if A.Data[i][i].Equal(A.Data[i][i].Zero()) {
			swapped := false
			for j := i + 1; j < n; j++ {
				if !A.Data[j][i].Equal(A.Data[j][i].Zero()) {
					A.Data[i], A.Data[j] = A.Data[j], A.Data[i]
					I.Data[i], I.Data[j] = I.Data[j], I.Data[i]
					swapped = true
					break
				}
			}
			if !swapped {
				return nil, errors.New("матрица вырождена")
			}
		}

		pivot := A.Data[i][i]
		for k := 0; k < n; k++ {
			f, _ := A.Data[i][k].Div(pivot)
			A.Data[i][k] = f
			f, _ = I.Data[i][k].Div(pivot)
			I.Data[i][k] = f
		}

		var wg sync.WaitGroup
		numCPU := runtime.NumCPU()
		ch := make(chan int, numCPU)

		for j := 0; j < n; j++ {
			if j == i {
				continue
			}
			ch <- 1
			wg.Add(1)
			go func(j int) {
				defer wg.Done()
				factor := A.Data[j][i]
				for k := 0; k < n; k++ {
					A.Data[j][k] = A.Data[j][k].Sub(factor.Mul(A.Data[i][k]))
					I.Data[j][k] = I.Data[j][k].Sub(factor.Mul(I.Data[i][k]))
				}
				<-ch
			}(j)
		}
		wg.Wait()
	}

	return I, nil
}
