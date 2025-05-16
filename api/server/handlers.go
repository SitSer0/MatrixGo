package server

import (
	"MatrixGo/internal/cache"
	"MatrixGo/internal/field"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var (
	// Создаем кэши для разных типов полей
	float64Cache  = cache.NewMatrixCache[field.Float64](5*time.Minute, 1000)
	complexCache  = cache.NewMatrixCache[field.Complex](5*time.Minute, 1000)
	rationalCache = cache.NewMatrixCache[field.Rational](5*time.Minute, 1000)
	gfCache       = cache.NewMatrixCache[field.GF](5*time.Minute, 1000)
)

func handleDeterminant(w http.ResponseWriter, r *http.Request) {
	var req MatrixRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "неверный формат запроса: "+err.Error(), http.StatusBadRequest)
		return
	}

	var result string
	var err error

	switch req.Type {
	case "float64":
		mat, err := ParseFloat64Matrix(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		det := float64Cache.CacheDeterminant(mat)
		result = det.String()

	case "complex":
		mat, err := ParseComplexMatrix(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		det := complexCache.CacheDeterminant(mat)
		result = det.String()

	case "rational":
		mat, err := ParseRationalMatrix(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		det := rationalCache.CacheDeterminant(mat)
		result = det.String()

	case "gf":
		mat, err := ParseGFMatrix(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		det := gfCache.CacheDeterminant(mat)
		result = det.String()

	default:
		http.Error(w, "неподдерживаемый тип поля", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"result": result})
}

func handleRank(w http.ResponseWriter, r *http.Request) {
	var req SystemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "неверный формат запроса: "+err.Error(), http.StatusBadRequest)
		return
	}

	var rank int
	var err error

	switch req.Matrix.Type {
	case "float64":
		mat, err := ParseFloat64Matrix(req.Matrix)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		rank = float64Cache.CacheRank(mat)

	case "complex":
		mat, err := ParseComplexMatrix(req.Matrix)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		rank = complexCache.CacheRank(mat)

	case "rational":
		mat, err := ParseRationalMatrix(req.Matrix)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		rank = rationalCache.CacheRank(mat)

	case "gf":
		mat, err := ParseGFMatrix(req.Matrix)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		rank = gfCache.CacheRank(mat)

	default:
		http.Error(w, "неподдерживаемый тип поля", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"value": rank})
}

func handleInverse(w http.ResponseWriter, r *http.Request) {
	var req MatrixRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "неверный формат запроса: "+err.Error(), http.StatusBadRequest)
		return
	}

	var err error
	var response [][]string

	switch req.Type {
	case "float64":
		mat, err := ParseFloat64Matrix(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		result, err := float64Cache.CacheInverse(mat)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		response = MatrixToStrings(result)

	case "complex":
		mat, err := ParseComplexMatrix(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		result, err := complexCache.CacheInverse(mat)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		response = MatrixToStrings(result)

	case "rational":
		mat, err := ParseRationalMatrix(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		result, err := rationalCache.CacheInverse(mat)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		response = MatrixToStrings(result)

	case "gf":
		mat, err := ParseGFMatrix(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		result, err := gfCache.CacheInverse(mat)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		response = MatrixToStrings(result)

	default:
		http.Error(w, "неподдерживаемый тип поля", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string][][]string{"result": response})
}

// HandleSolve обрабатывает запросы на решение системы уравнений
func HandleSolve(w http.ResponseWriter, r *http.Request) {
	// Добавляем CORS заголовки
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Читаем тело запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Ошибка чтения тела запроса: %v\n", err)
		http.Error(w, "ошибка чтения запроса", http.StatusBadRequest)
		return
	}
	fmt.Printf("Получено тело запроса: %s\n", string(body))

	// Восстанавливаем тело запроса для последующего использования
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	// Проверяем Content-Type
	contentType := r.Header.Get("Content-Type")
	fmt.Printf("Content-Type заголовок: %s\n", contentType)
	if !strings.Contains(contentType, "application/json") {
		fmt.Printf("Неверный Content-Type: %s\n", contentType)
		http.Error(w, "требуется application/json", http.StatusBadRequest)
		return
	}

	var req SystemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Printf("Ошибка декодирования JSON: %v\n", err)
		http.Error(w, fmt.Sprintf("неверный формат запроса: %v", err), http.StatusBadRequest)
		return
	}

	fmt.Printf("Получен запрос: %+v\n", req)
	fmt.Printf("Тип матрицы: %s\n", req.Matrix.Type)
	fmt.Printf("Размеры матрицы: %dx%d\n", req.Matrix.Rows, req.Matrix.Cols)
	fmt.Printf("Данные матрицы: %v\n", req.Matrix.Data)
	fmt.Printf("Длина вектора: %d\n", len(req.Vector))
	fmt.Printf("Данные вектора: %v\n", req.Vector)

	// Базовые проверки данных
	if req.Matrix.Rows <= 0 || req.Matrix.Cols <= 0 {
		errMsg := fmt.Sprintf("недопустимые размеры матрицы: %dx%d", req.Matrix.Rows, req.Matrix.Cols)
		fmt.Println(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	if len(req.Matrix.Data) != req.Matrix.Rows {
		errMsg := fmt.Sprintf("количество строк в данных (%d) не соответствует указанному размеру (%d)", len(req.Matrix.Data), req.Matrix.Rows)
		fmt.Println(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	for i, row := range req.Matrix.Data {
		if len(row) != req.Matrix.Cols {
			errMsg := fmt.Sprintf("количество столбцов в строке %d (%d) не соответствует указанному размеру (%d)", i, len(row), req.Matrix.Cols)
			fmt.Println(errMsg)
			http.Error(w, errMsg, http.StatusBadRequest)
			return
		}
	}

	if len(req.Vector) != req.Matrix.Rows {
		errMsg := fmt.Sprintf("размер вектора (%d) не соответствует количеству строк матрицы (%d)", len(req.Vector), req.Matrix.Rows)
		fmt.Println(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	var result []string
	var solveErr error

	switch req.Matrix.Type {
	case "float64":
		mat, err := ParseFloat64Matrix(req.Matrix)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		vec, err := ParseFloat64Vector(req.Vector)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		solution, err := float64Cache.CacheSolveSystem(mat, vec)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		result = VectorToStrings(solution)

	case "complex":
		mat, err := ParseComplexMatrix(req.Matrix)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		vec, err := ParseComplexVector(req.Vector)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		solution, err := complexCache.CacheSolveSystem(mat, vec)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		result = VectorToStrings(solution)

	case "rational":
		mat, err := ParseRationalMatrix(req.Matrix)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		vec, err := ParseRationalVector(req.Vector)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		solution, err := rationalCache.CacheSolveSystem(mat, vec)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		result = VectorToStrings(solution)

	case "gf":
		mat, err := ParseGFMatrix(req.Matrix)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		vec, err := ParseGFVector(req.Vector, req.Matrix.ModP)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		solution, err := gfCache.CacheSolveSystem(mat, vec)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		result = VectorToStrings(solution)

	default:
		http.Error(w, "неподдерживаемый тип поля", http.StatusBadRequest)
		return
	}

	if solveErr != nil {
		http.Error(w, solveErr.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string][]string{"result": result})
}
