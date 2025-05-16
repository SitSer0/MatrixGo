package server

import (
	"MatrixGo/internal/field"
	"MatrixGo/internal/matrix"
	"MatrixGo/internal/vector"
	"fmt"
	"strconv"
	"strings"
)

// ParseMatrix преобразует MatrixRequest в соответствующую матрицу нужного типа
func ParseMatrix(req MatrixRequest) (interface{}, error) {
	switch req.Type {
	case "float64":
		return ParseFloat64Matrix(req)
	case "complex":
		return ParseComplexMatrix(req)
	case "rational":
		return ParseRationalMatrix(req)
	case "gf":
		return ParseGFMatrix(req)
	default:
		return nil, fmt.Errorf("неподдерживаемый тип матрицы: %s", req.Type)
	}
}

func ParseFloat64Matrix(req MatrixRequest) (*matrix.Matrix[field.Float64], error) {
	if req.Rows <= 0 || req.Cols <= 0 {
		return nil, fmt.Errorf("недопустимые размеры матрицы: %dx%d", req.Rows, req.Cols)
	}

	if len(req.Data) != req.Rows {
		return nil, fmt.Errorf("количество строк в данных (%d) не соответствует указанному размеру (%d)", len(req.Data), req.Rows)
	}

	data := make([][]field.Float64, req.Rows)
	for i := range data {
		if len(req.Data[i]) != req.Cols {
			return nil, fmt.Errorf("количество столбцов в строке %d (%d) не соответствует указанному размеру (%d)", i, len(req.Data[i]), req.Cols)
		}
		data[i] = make([]field.Float64, req.Cols)
		for j := range data[i] {
			val, err := strconv.ParseFloat(req.Data[i][j], 64)
			if err != nil {
				return nil, fmt.Errorf("ошибка парсинга числа [%d][%d]: %w", i, j, err)
			}
			data[i][j] = field.Float64(val)
		}
	}
	return matrix.FromSlice(data)
}

func ParseComplexMatrix(req MatrixRequest) (*matrix.Matrix[field.Complex], error) {
	if req.Rows <= 0 || req.Cols <= 0 {
		return nil, fmt.Errorf("недопустимые размеры матрицы: %dx%d", req.Rows, req.Cols)
	}

	if len(req.Data) != req.Rows {
		return nil, fmt.Errorf("количество строк в данных (%d) не соответствует указанному размеру (%d)", len(req.Data), req.Rows)
	}

	data := make([][]field.Complex, req.Rows)
	for i := range data {
		if len(req.Data[i]) != req.Cols {
			return nil, fmt.Errorf("количество столбцов в строке %d (%d) не соответствует указанному размеру (%d)", i, len(req.Data[i]), req.Cols)
		}
		data[i] = make([]field.Complex, req.Cols)
		for j := range data[i] {
			val, err := field.ParseComplex(req.Data[i][j])
			if err != nil {
				return nil, fmt.Errorf("ошибка парсинга комплексного числа [%d][%d]: %w", i, j, err)
			}
			data[i][j] = val
		}
	}
	return matrix.FromSlice(data)
}

func ParseRationalMatrix(req MatrixRequest) (*matrix.Matrix[field.Rational], error) {
	data := make([][]field.Rational, req.Rows)
	for i := range data {
		data[i] = make([]field.Rational, req.Cols)
		for j := range data[i] {
			// Ожидаем формат "числитель/знаменатель" или просто "число"
			parts := strings.Split(req.Data[i][j], "/")
			var num, den int64
			var err error

			if len(parts) == 1 {
				num, err = strconv.ParseInt(parts[0], 10, 64)
				if err != nil {
					return nil, fmt.Errorf("ошибка парсинга числителя [%d][%d]: %w", i, j, err)
				}
				den = 1
			} else if len(parts) == 2 {
				num, err = strconv.ParseInt(parts[0], 10, 64)
				if err != nil {
					return nil, fmt.Errorf("ошибка парсинга числителя [%d][%d]: %w", i, j, err)
				}
				den, err = strconv.ParseInt(parts[1], 10, 64)
				if err != nil {
					return nil, fmt.Errorf("ошибка парсинга знаменателя [%d][%d]: %w", i, j, err)
				}
				if den == 0 {
					return nil, fmt.Errorf("знаменатель не может быть равен нулю [%d][%d]", i, j)
				}
			} else {
				return nil, fmt.Errorf("неверный формат рационального числа [%d][%d]", i, j)
			}

			data[i][j] = field.NewRational(num, den)
		}
	}
	return matrix.FromSlice(data)
}

func ParseGFMatrix(req MatrixRequest) (*matrix.Matrix[field.GF], error) {
	if req.ModP == 0 {
		return nil, fmt.Errorf("не указан модуль для конечного поля")
	}

	data := make([][]field.GF, req.Rows)
	for i := range data {
		data[i] = make([]field.GF, req.Cols)
		for j := range data[i] {
			val, err := strconv.ParseInt(req.Data[i][j], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("ошибка парсинга числа [%d][%d]: %w", i, j, err)
			}
			gf, err := field.NewGF(val, req.ModP)
			if err != nil {
				return nil, fmt.Errorf("ошибка создания элемента поля [%d][%d]: %w", i, j, err)
			}
			data[i][j] = gf
		}
	}
	return matrix.FromSlice(data)
}

// MatrixToStrings преобразует матрицу в двумерный массив строк для ответа
func MatrixToStrings(m interface{}) [][]string {
	switch mat := m.(type) {
	case *matrix.Matrix[field.Float64]:
		result := make([][]string, mat.Rows)
		for i := range result {
			result[i] = make([]string, mat.Cols)
			for j := range result[i] {
				result[i][j] = fmt.Sprintf("%v", mat.Data[i][j])
			}
		}
		return result
	case *matrix.Matrix[field.Complex]:
		result := make([][]string, mat.Rows)
		for i := range result {
			result[i] = make([]string, mat.Cols)
			for j := range result[i] {
				result[i][j] = fmt.Sprintf("%v", mat.Data[i][j])
			}
		}
		return result
	case *matrix.Matrix[field.Rational]:
		result := make([][]string, mat.Rows)
		for i := range result {
			result[i] = make([]string, mat.Cols)
			for j := range result[i] {
				result[i][j] = fmt.Sprintf("%v", mat.Data[i][j])
			}
		}
		return result
	case *matrix.Matrix[field.GF]:
		result := make([][]string, mat.Rows)
		for i := range result {
			result[i] = make([]string, mat.Cols)
			for j := range result[i] {
				result[i][j] = fmt.Sprintf("%v", mat.Data[i][j])
			}
		}
		return result
	default:
		return nil
	}
}

// ParseFloat64Vector конвертирует массив строк в вектор float64
func ParseFloat64Vector(data []string) (*vector.Vector[field.Float64], error) {
	if len(data) == 0 {
		fmt.Printf("Получен пустой вектор\n")
		return nil, fmt.Errorf("пустой вектор")
	}

	fmt.Printf("Парсинг вектора размера %d: %v\n", len(data), data)
	result := make([]field.Float64, len(data))
	for i, s := range data {
		if s == "" {
			fmt.Printf("Пустое значение в позиции %d\n", i)
			return nil, fmt.Errorf("пустое значение в позиции %d", i)
		}
		fmt.Printf("Парсинг числа в позиции %d: %q\n", i, s)
		val, err := strconv.ParseFloat(s, 64)
		if err != nil {
			fmt.Printf("Ошибка парсинга числа в позиции %d: %v\n", i, err)
			return nil, fmt.Errorf("ошибка парсинга числа в позиции %d: %w", i, err)
		}
		fmt.Printf("Успешно распарсено число в позиции %d: %v\n", i, val)
		result[i] = field.Float64(val)
	}
	return vector.NewVector(result), nil
}

// ParseComplexVector конвертирует массив строк в вектор комплексных чисел
func ParseComplexVector(data []string) (*vector.Vector[field.Complex], error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("пустой вектор")
	}

	result := make([]field.Complex, len(data))
	for i, s := range data {
		val, err := field.ParseComplex(s)
		if err != nil {
			return nil, fmt.Errorf("ошибка парсинга комплексного числа в позиции %d: %w", i, err)
		}
		result[i] = val
	}
	return vector.NewVector(result), nil
}

// ParseRationalVector конвертирует массив строк в вектор рациональных чисел
func ParseRationalVector(data []string) (*vector.Vector[field.Rational], error) {
	result := make([]field.Rational, len(data))
	for i, s := range data {
		parts := strings.Split(s, "/")
		var num, den int64
		var err error

		if len(parts) == 1 {
			num, err = strconv.ParseInt(parts[0], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("ошибка парсинга числителя в позиции %d: %w", i, err)
			}
			den = 1
		} else if len(parts) == 2 {
			num, err = strconv.ParseInt(parts[0], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("ошибка парсинга числителя в позиции %d: %w", i, err)
			}
			den, err = strconv.ParseInt(parts[1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("ошибка парсинга знаменателя в позиции %d: %w", i, err)
			}
			if den == 0 {
				return nil, fmt.Errorf("знаменатель не может быть равен нулю в позиции %d", i)
			}
		} else {
			return nil, fmt.Errorf("неверный формат рационального числа в позиции %d", i)
		}

		result[i] = field.NewRational(num, den)
	}
	return vector.NewVector(result), nil
}

// ParseGFVector конвертирует массив строк в вектор элементов конечного поля
func ParseGFVector(data []string, modP int64) (*vector.Vector[field.GF], error) {
	if modP == 0 {
		return nil, fmt.Errorf("не указан модуль для конечного поля")
	}

	result := make([]field.GF, len(data))
	for i, s := range data {
		val, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("ошибка парсинга числа в позиции %d: %w", i, err)
		}
		gf, err := field.NewGF(val, modP)
		if err != nil {
			return nil, fmt.Errorf("ошибка создания элемента поля в позиции %d: %w", i, err)
		}
		result[i] = gf
	}
	return vector.NewVector(result), nil
}

// VectorToStrings конвертирует вектор в массив строк для ответа
func VectorToStrings[T field.Field[T]](vec *vector.Vector[T]) []string {
	if vec == nil {
		return nil
	}
	result := make([]string, vec.Size)
	for i := range result {
		result[i] = fmt.Sprintf("%v", vec.Data[i])
	}
	return result
}
