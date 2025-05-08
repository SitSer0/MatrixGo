package matrix

import (
	"MatrixGo/internal/field"
	"MatrixGo/internal/vector"
	"errors"
)

func SolveSystem[T field.Field[T]](mat *Matrix[T], vec *vector.Vector[T]) (*vector.Vector[T], error) {
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
		if M.Data[i][i].Equal(M.Data[i][i].Zero()) {
			swapped := false
			for j := i + 1; j < n; j++ {
				if !M.Data[j][i].Equal(M.Data[j][i].Zero()) {
					M.Data[i], M.Data[j] = M.Data[j], M.Data[i]
					B[i], B[j] = B[j], B[i]
					swapped = true
					break
				}
			}
			if !swapped {
				return nil, errors.New("матрица вырождена, решение невозможно")
			}
		}

		pivot := M.Data[i][i]
		for j := i; j < n; j++ {
			d, _ := M.Data[i][j].Div(pivot)
			M.Data[i][j] = d
		}
		d, _ := B[i].Div(pivot)
		B[i] = d

		for k := i + 1; k < n; k++ {
			factor := M.Data[k][i]
			for j := i; j < n; j++ {
				M.Data[k][j] = M.Data[k][j].Sub(factor.Mul(M.Data[i][j]))
			}
			B[k] = B[k].Sub(factor.Mul(B[i]))
		}
	}

	x := make([]T, n)
	for i := n - 1; i >= 0; i-- {
		sum := B[i]
		for j := i + 1; j < n; j++ {
			sum = sum.Sub(M.Data[i][j].Mul(x[j]))
		}
		x[i] = sum
	}

	return vector.NewVector[T](x), nil
}

func SolveHomoSystem[T field.Field[T]](mat *Matrix[T]) (*vector.Vector[T], error) {
	vc := make([]T, mat.Rows)
	for i := 0; i < mat.Rows; i++ {
		vc[i] = mat.Data[0][0].Zero()
	}
	return SolveSystem(mat, vector.NewVector(vc))
}
