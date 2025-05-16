package server

import (
	"MatrixGo/internal/field"
	"MatrixGo/internal/matrix"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFloat64Matrix(t *testing.T) {
	tests := []struct {
		name    string
		req     MatrixRequest
		wantErr bool
	}{
		{
			name: "valid matrix",
			req: MatrixRequest{
				Type: "float64",
				Rows: 2,
				Cols: 2,
				Data: [][]string{
					{"1.5", "2.0"},
					{"3.0", "4.5"},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid number",
			req: MatrixRequest{
				Type: "float64",
				Rows: 2,
				Cols: 2,
				Data: [][]string{
					{"1.5", "invalid"},
					{"3.0", "4.5"},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mat, err := ParseFloat64Matrix(tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.req.Rows, mat.Rows)
			assert.Equal(t, tt.req.Cols, mat.Cols)
		})
	}
}

func TestParseComplexMatrix(t *testing.T) {
	tests := []struct {
		name    string
		req     MatrixRequest
		wantErr bool
	}{
		{
			name: "valid matrix",
			req: MatrixRequest{
				Type: "complex",
				Rows: 2,
				Cols: 2,
				Data: [][]string{
					{"1+2i", "3-4i"},
					{"5", "2.5i"},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid complex",
			req: MatrixRequest{
				Type: "complex",
				Rows: 2,
				Cols: 2,
				Data: [][]string{
					{"1+2i", "invalid"},
					{"5", "2.5i"},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mat, err := ParseComplexMatrix(tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.req.Rows, mat.Rows)
			assert.Equal(t, tt.req.Cols, mat.Cols)
		})
	}
}

func TestParseRationalMatrix(t *testing.T) {
	tests := []struct {
		name    string
		req     MatrixRequest
		wantErr bool
	}{
		{
			name: "valid matrix",
			req: MatrixRequest{
				Type: "rational",
				Rows: 2,
				Cols: 2,
				Data: [][]string{
					{"1/2", "3/4"},
					{"5", "6/1"},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid fraction",
			req: MatrixRequest{
				Type: "rational",
				Rows: 2,
				Cols: 2,
				Data: [][]string{
					{"1/2", "3/0"},
					{"5", "6/1"},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mat, err := ParseRationalMatrix(tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.req.Rows, mat.Rows)
			assert.Equal(t, tt.req.Cols, mat.Cols)
		})
	}
}

func TestParseGFMatrix(t *testing.T) {
	tests := []struct {
		name    string
		req     MatrixRequest
		wantErr bool
	}{
		{
			name: "valid matrix",
			req: MatrixRequest{
				Type: "gf",
				Rows: 2,
				Cols: 2,
				Data: [][]string{
					{"1", "2"},
					{"3", "4"},
				},
				ModP: 7,
			},
			wantErr: false,
		},
		{
			name: "invalid modulus",
			req: MatrixRequest{
				Type: "gf",
				Rows: 2,
				Cols: 2,
				Data: [][]string{
					{"1", "2"},
					{"3", "4"},
				},
				ModP: 4, // не простое число
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mat, err := ParseGFMatrix(tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.req.Rows, mat.Rows)
			assert.Equal(t, tt.req.Cols, mat.Cols)
		})
	}
}

func TestMatrixToStrings(t *testing.T) {
	t.Run("float64 matrix", func(t *testing.T) {
		data := [][]field.Float64{
			{field.Float64(1.5), field.Float64(2.0)},
			{field.Float64(3.0), field.Float64(4.5)},
		}
		mat, _ := matrix.FromSlice(data)
		result := MatrixToStrings(mat)
		assert.Equal(t, "1.5", result[0][0])
		assert.Equal(t, "2", result[0][1])
		assert.Equal(t, "3", result[1][0])
		assert.Equal(t, "4.5", result[1][1])
	})

	t.Run("nil matrix", func(t *testing.T) {
		result := MatrixToStrings(nil)
		assert.Nil(t, result)
	})
}
