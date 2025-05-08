package matrix

import (
	"MatrixGo/internal/field"
	"errors"
)

type Matrix[T field.Field[T]] struct {
	Rows int
	Cols int
	Data [][]T
}

func NewMatrix[T field.Field[T]](rows int, cols int, initVal T) *Matrix[T] {
	data := make([][]T, rows)

	for i := 0; i < rows; i++ {
		data[i] = make([]T, cols)
		for j := 0; j < cols; j++ {
			data[i][j] = initVal
		}
	}

	return &Matrix[T]{Rows: rows, Cols: cols, Data: data}
}

func (mat *Matrix[T]) Get(i, j int) (T, error) {
	if i < 0 || i >= mat.Cols || j < 0 || j >= mat.Rows {
		var zero T
		return zero, errors.New("индекс выходит за пределы матрицы")
	}

	return mat.Data[i][j], nil
}

func (mat *Matrix[T]) Set(i, j int, newVal T) error {
	if i < 0 || i >= mat.Cols || j < 0 || j >= mat.Rows {
		return errors.New("индекс выходит за пределы матрицы")
	}

	mat.Data[i][j] = newVal
	return nil
}

func (m1 *Matrix[T]) Add(m2 *Matrix[T]) (*Matrix[T], error) {
	if m1.Rows != m2.Rows || m1.Cols != m2.Cols {
		return nil, errors.New("матрицы не совпадают по размерам")
	}

	var zero T
	res := NewMatrix(m1.Rows, m1.Cols, zero)

	for i := 0; i < m1.Rows; i++ {
		for j := 0; j < m1.Cols; j++ {
			res.Data[i][j] = m1.Data[i][j].Add(m2.Data[i][j])
		}
	}

	return res, nil
}

func (m1 *Matrix[T]) Sub(m2 *Matrix[T]) (*Matrix[T], error) {
	if m1.Rows != m2.Rows || m1.Cols != m2.Cols {
		return nil, errors.New("матрицы не совпадают по размерам")
	}

	var zero T
	res := NewMatrix(m1.Rows, m1.Cols, zero)

	for i := 0; i < m1.Rows; i++ {
		for j := 0; j < m1.Cols; j++ {
			res.Data[i][j] = m1.Data[i][j].Sub(m2.Data[i][j])
		}
	}

	return res, nil
}

func (m1 *Matrix[T]) Mul(m2 *Matrix[T]) (*Matrix[T], error) {
	if m1.Cols != m2.Rows {
		return nil, errors.New("количество столбцов первой матрицы должно совпадать с количеством строк второй")
	}

	var zero T
	res := NewMatrix(m1.Rows, m2.Cols, zero)

	for i := 0; i < m1.Rows; i++ {
		for j := 0; j < m2.Cols; j++ {
			elem := m1.Data[i][0].Mul(m2.Data[0][j])
			for k := 1; k < m1.Cols; k++ {
				elem = elem.Add(m1.Data[i][k].Mul(m2.Data[k][j]))
			}
			res.Data[i][j] = elem
		}
	}

	return res, nil
}

func (m *Matrix[T]) Transpose() *Matrix[T] {
	var zero T
	res := NewMatrix(m.Cols, m.Rows, zero)
	for i := 0; i < m.Rows; i++ {
		for j := 0; j < m.Cols; j++ {
			res.Data[j][i] = m.Data[i][j]
		}
	}
	return res
}

func (m *Matrix[T]) Trace() T {
	sum := m.Data[0][0]
	for i := 1; i < m.Rows && i < m.Cols; i++ {
		sum = sum.Add(m.Data[i][i])
	}
	return sum
}

func IdentityMatrix[T field.Field[T]](size int, zero, one T) *Matrix[T] {
	identity := NewMatrix[T](size, size, zero)
	for i := 0; i < size; i++ {
		identity.Data[i][i] = one
	}
	return identity
}

func (m *Matrix[T]) Clone() *Matrix[T] {
	cloned := NewMatrix[T](m.Rows, m.Cols, m.Data[0][0].Zero())
	for i := 0; i < m.Rows; i++ {
		for j := 0; j < m.Cols; j++ {
			cloned.Data[i][j] = m.Data[i][j]
		}
	}
	return cloned
}
