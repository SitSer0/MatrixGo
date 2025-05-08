package matrix

import (
	"MatrixGo/internal/field"
	"errors"
	"runtime"
	"sync"
)

// Task описывает координаты одного элемента, который нужно вычислить
type multiplyTask[T field.Field[T]] struct {
	Row int
	Col int
}

// MultiplyParallel выполняет параллельное умножение матриц
func (m1 *Matrix[T]) MultiplyParallel(m2 *Matrix[T]) (*Matrix[T], error) {
	if m1.Cols != m2.Rows {
		return nil, errors.New("недопустимые размеры для умножения")
	}

	res := NewMatrix[T](m1.Rows, m2.Cols, m1.Data[0][0].Zero())

	numWorkers := runtime.NumCPU()
	tasks := make(chan multiplyTask[T], m1.Rows*m2.Cols)
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range tasks {
				sum := m1.Data[0][0].Zero()
				for k := 0; k < m1.Cols; k++ {
					sum = sum.Add(m1.Data[task.Row][k].Mul(m2.Data[k][task.Col]))
				}
				res.Data[task.Row][task.Col] = sum
			}
		}()
	}

	for i := 0; i < m1.Rows; i++ {
		for j := 0; j < m2.Cols; j++ {
			tasks <- multiplyTask[T]{Row: i, Col: j}
		}
	}
	close(tasks)

	wg.Wait()

	return res, nil
}
