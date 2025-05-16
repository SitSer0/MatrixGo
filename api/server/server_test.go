package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServer_HandleMatrixAdd(t *testing.T) {
	s := NewServer()

	tests := []struct {
		name       string
		request    map[string]interface{}
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid float64 matrices",
			request: map[string]interface{}{
				"matrix1": MatrixRequest{
					Type: "float64",
					Rows: 2,
					Cols: 2,
					Data: [][]string{
						{"1.0", "2.0"},
						{"3.0", "4.0"},
					},
				},
				"matrix2": MatrixRequest{
					Type: "float64",
					Rows: 2,
					Cols: 2,
					Data: [][]string{
						{"5.0", "6.0"},
						{"7.0", "8.0"},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "different types",
			request: map[string]interface{}{
				"matrix1": MatrixRequest{
					Type: "float64",
					Rows: 2,
					Cols: 2,
					Data: [][]string{
						{"1.0", "2.0"},
						{"3.0", "4.0"},
					},
				},
				"matrix2": MatrixRequest{
					Type: "complex",
					Rows: 2,
					Cols: 2,
					Data: [][]string{
						{"5.0", "6.0"},
						{"7.0", "8.0"},
					},
				},
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.request)
			req := httptest.NewRequest("POST", "/api/v1/matrix/add", bytes.NewReader(body))
			w := httptest.NewRecorder()

			s.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			if !tt.wantErr {
				var response MatrixResponse
				err := json.NewDecoder(w.Body).Decode(&response)
				assert.NoError(t, err)
				assert.NotNil(t, response.Result)
				assert.Empty(t, response.Error)
			}
		})
	}
}

func TestServer_HandleMatrixMultiply(t *testing.T) {
	s := NewServer()

	tests := []struct {
		name       string
		request    map[string]interface{}
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid rational matrices",
			request: map[string]interface{}{
				"matrix1": MatrixRequest{
					Type: "rational",
					Rows: 2,
					Cols: 2,
					Data: [][]string{
						{"1/2", "2/3"},
						{"3/4", "4/5"},
					},
				},
				"matrix2": MatrixRequest{
					Type: "rational",
					Rows: 2,
					Cols: 2,
					Data: [][]string{
						{"1/3", "1/4"},
						{"1/2", "2/3"},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "incompatible dimensions",
			request: map[string]interface{}{
				"matrix1": MatrixRequest{
					Type: "float64",
					Rows: 2,
					Cols: 3,
					Data: [][]string{
						{"1.0", "2.0", "3.0"},
						{"4.0", "5.0", "6.0"},
					},
				},
				"matrix2": MatrixRequest{
					Type: "float64",
					Rows: 2,
					Cols: 2,
					Data: [][]string{
						{"1.0", "2.0"},
						{"3.0", "4.0"},
					},
				},
			},
			wantStatus: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.request)
			req := httptest.NewRequest("POST", "/api/v1/matrix/multiply", bytes.NewReader(body))
			w := httptest.NewRecorder()

			s.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			if !tt.wantErr {
				var response MatrixResponse
				err := json.NewDecoder(w.Body).Decode(&response)
				assert.NoError(t, err)
				assert.NotNil(t, response.Result)
				assert.Empty(t, response.Error)
			}
		})
	}
}

func TestServer_HandleMatrixTranspose(t *testing.T) {
	s := NewServer()

	tests := []struct {
		name       string
		request    MatrixRequest
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid gf matrix",
			request: MatrixRequest{
				Type: "gf",
				Rows: 2,
				Cols: 2,
				Data: [][]string{
					{"1", "2"},
					{"3", "4"},
				},
				ModP: 7,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid matrix type",
			request: MatrixRequest{
				Type: "unknown",
				Rows: 2,
				Cols: 2,
				Data: [][]string{
					{"1", "2"},
					{"3", "4"},
				},
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.request)
			req := httptest.NewRequest("POST", "/api/v1/matrix/transpose", bytes.NewReader(body))
			w := httptest.NewRecorder()

			s.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			if !tt.wantErr {
				var response MatrixResponse
				err := json.NewDecoder(w.Body).Decode(&response)
				assert.NoError(t, err)
				assert.NotNil(t, response.Result)
				assert.Empty(t, response.Error)
			}
		})
	}
}

func TestServer_InvalidJSON(t *testing.T) {
	s := NewServer()

	req := httptest.NewRequest("POST", "/api/v1/matrix/add", bytes.NewReader([]byte("invalid json")))
	w := httptest.NewRecorder()

	s.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
