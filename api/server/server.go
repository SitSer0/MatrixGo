package server

import (
	"MatrixGo/internal/field"
	"MatrixGo/internal/matrix"
	"MatrixGo/internal/vector"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type Server struct {
	router *mux.Router
}

func NewServer() *Server {
	s := &Server{
		router: mux.NewRouter(),
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.router.HandleFunc("/api/v1/matrix/add", s.handleMatrixAdd()).Methods("POST")
	s.router.HandleFunc("/api/v1/matrix/multiply", s.handleMatrixMultiply()).Methods("POST")
	s.router.HandleFunc("/api/v1/matrix/transpose", s.handleMatrixTranspose()).Methods("POST")
	s.router.HandleFunc("/api/v1/matrix/solve", s.handleMatrixSolve()).Methods("POST")
	s.router.HandleFunc("/api/v1/matrix/inverse", s.handleMatrixInverse()).Methods("POST")
	s.router.HandleFunc("/api/v1/matrix/determinant", s.handleMatrixDeterminant()).Methods("POST")
	s.router.HandleFunc("/api/v1/matrix/rank", s.handleMatrixRank()).Methods("POST")
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) handleMatrixAdd() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Matrix1 MatrixRequest `json:"matrix1"`
			Matrix2 MatrixRequest `json:"matrix2"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Проверяем, что матрицы одного типа
		if req.Matrix1.Type != req.Matrix2.Type {
			http.Error(w, "матрицы разных типов", http.StatusBadRequest)
			return
		}

		// Парсим матрицы
		m1, err := ParseMatrix(req.Matrix1)
		if err != nil {
			http.Error(w, fmt.Sprintf("ошибка парсинга первой матрицы: %v", err), http.StatusBadRequest)
			return
		}

		m2, err := ParseMatrix(req.Matrix2)
		if err != nil {
			http.Error(w, fmt.Sprintf("ошибка парсинга второй матрицы: %v", err), http.StatusBadRequest)
			return
		}

		// Выполняем сложение в зависимости от типа
		var result interface{}
		switch m1 := m1.(type) {
		case *matrix.Matrix[field.Float64]:
			m2 := m2.(*matrix.Matrix[field.Float64])
			result, err = m1.Add(m2)
		case *matrix.Matrix[field.Complex]:
			m2 := m2.(*matrix.Matrix[field.Complex])
			result, err = m1.Add(m2)
		case *matrix.Matrix[field.Rational]:
			m2 := m2.(*matrix.Matrix[field.Rational])
			result, err = m1.Add(m2)
		case *matrix.Matrix[field.GF]:
			m2 := m2.(*matrix.Matrix[field.GF])
			result, err = m1.Add(m2)
		default:
			http.Error(w, "неподдерживаемый тип матрицы", http.StatusBadRequest)
			return
		}

		if err != nil {
			http.Error(w, fmt.Sprintf("ошибка сложения матриц: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(MatrixResponse{
			Result: MatrixToStrings(result),
		})
	}
}

func (s *Server) handleMatrixMultiply() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Matrix1 MatrixRequest `json:"matrix1"`
			Matrix2 MatrixRequest `json:"matrix2"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Проверяем, что матрицы одного типа
		if req.Matrix1.Type != req.Matrix2.Type {
			http.Error(w, "матрицы разных типов", http.StatusBadRequest)
			return
		}

		// Парсим матрицы
		m1, err := ParseMatrix(req.Matrix1)
		if err != nil {
			http.Error(w, fmt.Sprintf("ошибка парсинга первой матрицы: %v", err), http.StatusBadRequest)
			return
		}

		m2, err := ParseMatrix(req.Matrix2)
		if err != nil {
			http.Error(w, fmt.Sprintf("ошибка парсинга второй матрицы: %v", err), http.StatusBadRequest)
			return
		}

		// Выполняем умножение в зависимости от типа
		var result interface{}
		switch m1 := m1.(type) {
		case *matrix.Matrix[field.Float64]:
			m2 := m2.(*matrix.Matrix[field.Float64])
			result, err = m1.Mul(m2)
		case *matrix.Matrix[field.Complex]:
			m2 := m2.(*matrix.Matrix[field.Complex])
			result, err = m1.Mul(m2)
		case *matrix.Matrix[field.Rational]:
			m2 := m2.(*matrix.Matrix[field.Rational])
			result, err = m1.Mul(m2)
		case *matrix.Matrix[field.GF]:
			m2 := m2.(*matrix.Matrix[field.GF])
			result, err = m1.Mul(m2)
		default:
			http.Error(w, "неподдерживаемый тип матрицы", http.StatusBadRequest)
			return
		}

		if err != nil {
			http.Error(w, fmt.Sprintf("ошибка умножения матриц: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(MatrixResponse{
			Result: MatrixToStrings(result),
		})
	}
}

func (s *Server) handleMatrixTranspose() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req MatrixRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Парсим матрицу
		m, err := ParseMatrix(req)
		if err != nil {
			http.Error(w, fmt.Sprintf("ошибка парсинга матрицы: %v", err), http.StatusBadRequest)
			return
		}

		// Выполняем транспонирование в зависимости от типа
		var result interface{}
		switch m := m.(type) {
		case *matrix.Matrix[field.Float64]:
			result = m.Transpose()
		case *matrix.Matrix[field.Complex]:
			result = m.Transpose()
		case *matrix.Matrix[field.Rational]:
			result = m.Transpose()
		case *matrix.Matrix[field.GF]:
			result = m.Transpose()
		default:
			http.Error(w, "неподдерживаемый тип матрицы", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(MatrixResponse{
			Result: MatrixToStrings(result),
		})
	}
}

func (s *Server) handleMatrixSolve() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SystemRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Парсим матрицу
		m, err := ParseMatrix(req.Matrix)
		if err != nil {
			http.Error(w, fmt.Sprintf("ошибка парсинга матрицы: %v", err), http.StatusBadRequest)
			return
		}

		// Парсим вектор правой части
		var b interface{}
		switch req.Matrix.Type {
		case "float64":
			data := make([]field.Float64, len(req.Vector))
			for i, v := range req.Vector {
				val, err := strconv.ParseFloat(v, 64)
				if err != nil {
					http.Error(w, fmt.Sprintf("ошибка парсинга вектора [%d]: %v", i, err), http.StatusBadRequest)
					return
				}
				data[i] = field.Float64(val)
			}
			b = vector.NewVector(data)
		case "complex":
			data := make([]field.Complex, len(req.Vector))
			for i, v := range req.Vector {
				val, err := field.ParseComplex(v)
				if err != nil {
					http.Error(w, fmt.Sprintf("ошибка парсинга вектора [%d]: %v", i, err), http.StatusBadRequest)
					return
				}
				data[i] = val
			}
			b = vector.NewVector(data)
		case "rational":
			data := make([]field.Rational, len(req.Vector))
			for i, v := range req.Vector {
				parts := strings.Split(v, "/")
				var num, den int64
				if len(parts) == 1 {
					num, err = strconv.ParseInt(parts[0], 10, 64)
					if err != nil {
						http.Error(w, fmt.Sprintf("ошибка парсинга вектора [%d]: %v", i, err), http.StatusBadRequest)
						return
					}
					den = 1
				} else if len(parts) == 2 {
					num, err = strconv.ParseInt(parts[0], 10, 64)
					if err != nil {
						http.Error(w, fmt.Sprintf("ошибка парсинга вектора [%d]: %v", i, err), http.StatusBadRequest)
						return
					}
					den, err = strconv.ParseInt(parts[1], 10, 64)
					if err != nil {
						http.Error(w, fmt.Sprintf("ошибка парсинга вектора [%d]: %v", i, err), http.StatusBadRequest)
						return
					}
					if den == 0 {
						http.Error(w, fmt.Sprintf("знаменатель не может быть равен нулю [%d]", i), http.StatusBadRequest)
						return
					}
				} else {
					http.Error(w, fmt.Sprintf("неверный формат рационального числа [%d]", i), http.StatusBadRequest)
					return
				}
				data[i] = field.NewRational(num, den)
			}
			b = vector.NewVector(data)
		case "gf":
			data := make([]field.GF, len(req.Vector))
			for i, v := range req.Vector {
				val, err := strconv.ParseInt(v, 10, 64)
				if err != nil {
					http.Error(w, fmt.Sprintf("ошибка парсинга вектора [%d]: %v", i, err), http.StatusBadRequest)
					return
				}
				gf, err := field.NewGF(val, req.Matrix.ModP)
				if err != nil {
					http.Error(w, fmt.Sprintf("ошибка создания элемента поля [%d]: %v", i, err), http.StatusBadRequest)
					return
				}
				data[i] = gf
			}
			b = vector.NewVector(data)
		default:
			http.Error(w, "неподдерживаемый тип матрицы", http.StatusBadRequest)
			return
		}

		// Решаем систему в зависимости от типа
		var result interface{}
		switch m := m.(type) {
		case *matrix.Matrix[field.Float64]:
			b := b.(*vector.Vector[field.Float64])
			result, err = matrix.SolveSystem(m, b)
		case *matrix.Matrix[field.Complex]:
			b := b.(*vector.Vector[field.Complex])
			result, err = matrix.SolveSystem(m, b)
		case *matrix.Matrix[field.Rational]:
			b := b.(*vector.Vector[field.Rational])
			result, err = matrix.SolveSystem(m, b)
		case *matrix.Matrix[field.GF]:
			b := b.(*vector.Vector[field.GF])
			result, err = matrix.SolveSystem(m, b)
		default:
			http.Error(w, "неподдерживаемый тип матрицы", http.StatusBadRequest)
			return
		}

		if err != nil {
			http.Error(w, fmt.Sprintf("ошибка решения системы: %v", err), http.StatusInternalServerError)
			return
		}

		// Преобразуем результат в строки
		var resultStrings []string
		switch result := result.(type) {
		case *vector.Vector[field.Float64]:
			resultStrings = make([]string, result.Size)
			for i := 0; i < result.Size; i++ {
				resultStrings[i] = fmt.Sprintf("%v", result.Data[i])
			}
		case *vector.Vector[field.Complex]:
			resultStrings = make([]string, result.Size)
			for i := 0; i < result.Size; i++ {
				resultStrings[i] = fmt.Sprintf("%v", result.Data[i])
			}
		case *vector.Vector[field.Rational]:
			resultStrings = make([]string, result.Size)
			for i := 0; i < result.Size; i++ {
				resultStrings[i] = fmt.Sprintf("%v", result.Data[i])
			}
		case *vector.Vector[field.GF]:
			resultStrings = make([]string, result.Size)
			for i := 0; i < result.Size; i++ {
				resultStrings[i] = fmt.Sprintf("%v", result.Data[i])
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(MatrixResponse{
			Result: [][]string{resultStrings}, // Возвращаем как матрицу 1xn
		})
	}
}

func (s *Server) handleMatrixInverse() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req MatrixRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Парсим матрицу
		m, err := ParseMatrix(req)
		if err != nil {
			http.Error(w, fmt.Sprintf("ошибка парсинга матрицы: %v", err), http.StatusBadRequest)
			return
		}

		// Выполняем обращение в зависимости от типа
		var result interface{}
		switch m := m.(type) {
		case *matrix.Matrix[field.Float64]:
			result, err = m.Inverse()
		case *matrix.Matrix[field.Complex]:
			result, err = m.Inverse()
		case *matrix.Matrix[field.Rational]:
			result, err = m.Inverse()
		case *matrix.Matrix[field.GF]:
			result, err = m.Inverse()
		default:
			http.Error(w, "неподдерживаемый тип матрицы", http.StatusBadRequest)
			return
		}

		if err != nil {
			http.Error(w, fmt.Sprintf("ошибка обращения матрицы: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(MatrixResponse{
			Result: MatrixToStrings(result),
		})
	}
}

func (s *Server) handleMatrixDeterminant() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req MatrixRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Парсим матрицу
		m, err := ParseMatrix(req)
		if err != nil {
			http.Error(w, fmt.Sprintf("ошибка парсинга матрицы: %v", err), http.StatusBadRequest)
			return
		}

		// Выполняем вычисление определителя в зависимости от типа
		var result interface{}
		switch m := m.(type) {
		case *matrix.Matrix[field.Float64]:
			result = m.Determinant()
		case *matrix.Matrix[field.Complex]:
			result = m.Determinant()
		case *matrix.Matrix[field.Rational]:
			result = m.Determinant()
		case *matrix.Matrix[field.GF]:
			result = m.Determinant()
		default:
			http.Error(w, "неподдерживаемый тип матрицы", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(MatrixResponse{
			Result: MatrixToStrings(result),
		})
	}
}

func (s *Server) handleMatrixRank() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req MatrixRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Парсим матрицу
		m, err := ParseMatrix(req)
		if err != nil {
			http.Error(w, fmt.Sprintf("ошибка парсинга матрицы: %v", err), http.StatusBadRequest)
			return
		}

		// Выполняем вычисление ранга в зависимости от типа
		var result interface{}
		switch m := m.(type) {
		case *matrix.Matrix[field.Float64]:
			result = m.Rank()
		case *matrix.Matrix[field.Complex]:
			result = m.Rank()
		case *matrix.Matrix[field.Rational]:
			result = m.Rank()
		case *matrix.Matrix[field.GF]:
			result = m.Rank()
		default:
			http.Error(w, "неподдерживаемый тип матрицы", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(MatrixResponse{
			Result: MatrixToStrings(result),
		})
	}
}
