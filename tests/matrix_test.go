package test

import (
	"MatrixGo/internal/field"
	"MatrixGo/internal/matrix"
	"MatrixGo/internal/vector"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMatrixInitialization(t *testing.T) {
	mat := matrix.NewMatrix[field.Float64](2, 2, field.Float64(0))
	assert.NotNil(t, mat)
	assert.Equal(t, 2, mat.Rows)
	assert.Equal(t, 2, mat.Cols)
	assert.Equal(t, field.Float64(0), mat.Data[0][0])
}

func TestMatrixAdd(t *testing.T) {
	mat1 := matrix.NewMatrix[field.Float64](2, 2, field.Float64(1))
	mat2 := matrix.NewMatrix[field.Float64](2, 2, field.Float64(2))

	result, err := mat1.Add(mat2)
	assert.NoError(t, err)
	assert.Equal(t, field.Float64(3), result.Data[0][0])
	assert.Equal(t, field.Float64(3), result.Data[1][1])
}

func TestMatrixAddError(t *testing.T) {
	mat1 := matrix.NewMatrix[field.Float64](2, 2, field.Float64(1))
	mat2 := matrix.NewMatrix[field.Float64](3, 3, field.Float64(2))

	_, err := mat1.Add(mat2)
	assert.Error(t, err)
}

func TestVectorInitialization(t *testing.T) {
	vec := vector.NewVector([]field.Float64{field.Float64(1), field.Float64(2)})
	assert.NotNil(t, vec)
	assert.Equal(t, 2, vec.Size)
	assert.Equal(t, field.Float64(1), vec.Data[0])
}

func TestSolveSystem(t *testing.T) {
	A := matrix.NewMatrix[field.Float64](3, 3, field.Float64(0))
	A.Data[0] = []field.Float64{field.Float64(3), field.Float64(2), field.Float64(-1)}
	A.Data[1] = []field.Float64{field.Float64(2), field.Float64(-1), field.Float64(3)}
	A.Data[2] = []field.Float64{-1, field.Float64(3), field.Float64(2)}

	b := vector.NewVector([]field.Float64{
		field.Float64(1),
		field.Float64(2),
		field.Float64(3),
	})

	result, err := matrix.SolveSystem(A, b)
	assert.NoError(t, err)

	assert.Equal(t, field.Float64(1), result.Data[0])
	assert.Equal(t, field.Float64(1), result.Data[1])
	assert.Equal(t, field.Float64(1), result.Data[2])
}

func TestSolveHomoSystem(t *testing.T) {

	A := matrix.NewMatrix[field.Float64](3, 3, field.Float64(0))
	A.Data[0] = []field.Float64{field.Float64(1), field.Float64(2), field.Float64(3)}
	A.Data[1] = []field.Float64{field.Float64(4), field.Float64(5), field.Float64(6)}
	A.Data[2] = []field.Float64{field.Float64(7), field.Float64(8), field.Float64(9)}

	result, err := matrix.SolveHomoSystem(A)
	assert.NoError(t, err)

	assert.Equal(t, field.Float64(0), result.Data[0])
	assert.Equal(t, field.Float64(0), result.Data[1])
	assert.Equal(t, field.Float64(0), result.Data[2])
}

func TestMatrixTranspose(t *testing.T) {
	mat := matrix.NewMatrix[field.Float64](2, 3, field.Float64(0))
	mat.Data[0] = []field.Float64{field.Float64(1), field.Float64(2), field.Float64(3)}
	mat.Data[1] = []field.Float64{field.Float64(4), field.Float64(5), field.Float64(6)}

	transposed := mat.Transpose()

	assert.Equal(t, field.Float64(1), transposed.Data[0][0])
	assert.Equal(t, field.Float64(4), transposed.Data[0][1])
	assert.Equal(t, field.Float64(2), transposed.Data[1][0])
	assert.Equal(t, field.Float64(5), transposed.Data[1][1])
}
