package matrix

import (
	"MatrixGo/internal/field"
	"MatrixGo/internal/vector"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDeterminantParallel(t *testing.T) {
	tests := []struct {
		name     string
		matrix   [][]field.Float64
		expected field.Float64
	}{
		{
			name: "2x2 matrix",
			matrix: [][]field.Float64{
				{1, 2},
				{3, 4},
			},
			expected: -2,
		},
		{
			name: "3x3 matrix",
			matrix: [][]field.Float64{
				{1, 2, 3},
				{4, 5, 6},
				{7, 8, 9},
			},
			expected: 0,
		},
		{
			name: "4x4 matrix",
			matrix: [][]field.Float64{
				{1, 2, 3, 4},
				{5, 6, 7, 8},
				{9, 10, 11, 12},
				{13, 14, 15, 16},
			},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mat, _ := FromSlice(tt.matrix)

			// Замеряем время выполнения обычного и параллельного алгоритмов
			start := time.Now()
			det1 := mat.Determinant()
			serial := time.Since(start)

			start = time.Now()
			det2 := mat.DeterminantParallel()
			parallel := time.Since(start)

			// Проверяем корректность результата
			assert.Equal(t, tt.expected, det1)
			assert.Equal(t, tt.expected, det2)

			// Для больших матриц параллельная версия должна быть быстрее
			if len(tt.matrix) > 100 {
				assert.Less(t, parallel, serial)
			}
		})
	}
}

func TestRankParallel(t *testing.T) {
	tests := []struct {
		name     string
		matrix   [][]field.Float64
		expected int
	}{
		{
			name: "full rank 2x2",
			matrix: [][]field.Float64{
				{1, 0},
				{0, 1},
			},
			expected: 2,
		},
		{
			name: "rank 2 3x3",
			matrix: [][]field.Float64{
				{1, 2, 3},
				{2, 4, 6},
				{3, 6, 9},
			},
			expected: 1,
		},
		{
			name: "zero matrix",
			matrix: [][]field.Float64{
				{0, 0},
				{0, 0},
			},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mat, _ := FromSlice(tt.matrix)

			// Замеряем время выполнения обычного и параллельного алгоритмов
			start := time.Now()
			rank1 := mat.Rank()
			serial := time.Since(start)

			start = time.Now()
			rank2 := mat.RankParallel()
			parallel := time.Since(start)

			// Проверяем корректность результата
			assert.Equal(t, tt.expected, rank1)
			assert.Equal(t, tt.expected, rank2)

			// Для больших матриц параллельная версия должна быть быстрее
			if len(tt.matrix) > 100 {
				assert.Less(t, parallel, serial)
			}
		})
	}
}

func TestSolveSystemParallel(t *testing.T) {
	tests := []struct {
		name     string
		matrix   [][]field.Float64
		vector   []field.Float64
		expected []field.Float64
		wantErr  bool
	}{
		{
			name: "simple system",
			matrix: [][]field.Float64{
				{2, 1},
				{1, 3},
			},
			vector:   []field.Float64{4, 5},
			expected: []field.Float64{7.0 / 5.0, 6.0 / 5.0},
			wantErr:  false,
		},
		{
			name: "singular matrix",
			matrix: [][]field.Float64{
				{1, 1},
				{1, 1},
			},
			vector:   []field.Float64{1, 1},
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mat, _ := FromSlice(tt.matrix)
			vec := vector.NewVector(tt.vector)

			// Замеряем время выполнения обычного и параллельного алгоритмов
			start := time.Now()
			sol1, err1 := SolveSystem(mat, vec)
			serial := time.Since(start)

			start = time.Now()
			sol2, err2 := SolveSystemParallel(mat, vec)
			parallel := time.Since(start)

			if tt.wantErr {
				assert.Error(t, err1)
				assert.Error(t, err2)
				return
			}

			assert.NoError(t, err1)
			assert.NoError(t, err2)

			// Проверяем корректность результата
			for i := 0; i < len(tt.expected); i++ {
				assert.InDelta(t, float64(tt.expected[i]), float64(sol1.Data[i]), 1e-10)
				assert.InDelta(t, float64(tt.expected[i]), float64(sol2.Data[i]), 1e-10)
			}

			// Для больших матриц параллельная версия должна быть быстрее
			if len(tt.matrix) > 100 {
				assert.Less(t, parallel, serial)
			}
		})
	}
}

func TestInverseParallel(t *testing.T) {
	tests := []struct {
		name     string
		matrix   [][]field.Float64
		expected [][]field.Float64
		wantErr  bool
	}{
		{
			name: "2x2 matrix",
			matrix: [][]field.Float64{
				{4, 7},
				{2, 6},
			},
			expected: [][]field.Float64{
				{0.6, -0.7},
				{-0.2, 0.4},
			},
			wantErr: false,
		},
		{
			name: "singular matrix",
			matrix: [][]field.Float64{
				{1, 1},
				{1, 1},
			},
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mat, _ := FromSlice(tt.matrix)

			// Замеряем время выполнения обычного и параллельного алгоритмов
			start := time.Now()
			inv1, err1 := mat.Inverse()
			serial := time.Since(start)

			start = time.Now()
			inv2, err2 := mat.InverseParallel()
			parallel := time.Since(start)

			if tt.wantErr {
				assert.Error(t, err1)
				assert.Error(t, err2)
				return
			}

			assert.NoError(t, err1)
			assert.NoError(t, err2)

			// Проверяем корректность результата
			for i := 0; i < len(tt.expected); i++ {
				for j := 0; j < len(tt.expected[i]); j++ {
					assert.InDelta(t, float64(tt.expected[i][j]), float64(inv1.Data[i][j]), 1e-10)
					assert.InDelta(t, float64(tt.expected[i][j]), float64(inv2.Data[i][j]), 1e-10)
				}
			}

			// Для больших матриц параллельная версия должна быть быстрее
			if len(tt.matrix) > 100 {
				assert.Less(t, parallel, serial)
			}
		})
	}
}
