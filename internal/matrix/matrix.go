package matrix

import (
	"MatrixGo/internal/field"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
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

// String возвращает строковое представление матрицы в удобочитаемом формате
func (mat *Matrix[T]) String() string {
	if mat == nil || len(mat.Data) == 0 {
		return "[]"
	}

	// Находим максимальную длину строкового представления элементов
	maxLen := 0
	for i := 0; i < mat.Rows; i++ {
		for j := 0; j < mat.Cols; j++ {
			strLen := len(fmt.Sprintf("%v", mat.Data[i][j]))
			if strLen > maxLen {
				maxLen = strLen
			}
		}
	}

	// Формируем строку с выравниванием
	var sb strings.Builder
	sb.WriteString("⎡")
	for i := 0; i < mat.Rows; i++ {
		if i > 0 {
			sb.WriteString("⎢")
		}
		for j := 0; j < mat.Cols; j++ {
			if j == 0 {
				sb.WriteString(" ")
			}
			format := fmt.Sprintf("%%%dv", maxLen)
			sb.WriteString(fmt.Sprintf(format, mat.Data[i][j]))
			if j < mat.Cols-1 {
				sb.WriteString("  ")
			}
		}
		if i < mat.Rows-1 {
			sb.WriteString(" ⎥\n")
		} else {
			sb.WriteString(" ⎦")
		}
	}
	return sb.String()
}

// Zeros создает матрицу, заполненную нулями
func Zeros[T field.Field[T]](rows, cols int, zero T) *Matrix[T] {
	return NewMatrix(rows, cols, zero)
}

// Ones создает матрицу, заполненную единицами
func Ones[T field.Field[T]](rows, cols int, one T) *Matrix[T] {
	return NewMatrix(rows, cols, one)
}

// Eye создает единичную матрицу
func Eye[T field.Field[T]](size int, zero, one T) *Matrix[T] {
	mat := Zeros(size, size, zero)
	for i := 0; i < size; i++ {
		mat.Data[i][i] = one
	}
	return mat
}

// FromSlice создает матрицу из двумерного среза
func FromSlice[T field.Field[T]](data [][]T) (*Matrix[T], error) {
	if len(data) == 0 || len(data[0]) == 0 {
		return nil, errors.New("пустые данные")
	}

	rows := len(data)
	cols := len(data[0])

	// Проверяем, что все строки имеют одинаковую длину
	for i := 1; i < rows; i++ {
		if len(data[i]) != cols {
			return nil, errors.New("неравномерные данные")
		}
	}

	mat := Zeros(rows, cols, data[0][0].Zero())
	for i := 0; i < rows; i++ {
		copy(mat.Data[i], data[i])
	}
	return mat, nil
}

// MarshalJSON реализует интерфейс json.Marshaler
func (mat *Matrix[T]) MarshalJSON() ([]byte, error) {
	type MatrixJSON struct {
		Rows int        `json:"rows"`
		Cols int        `json:"cols"`
		Data [][]string `json:"data"`
	}

	jsonData := MatrixJSON{
		Rows: mat.Rows,
		Cols: mat.Cols,
		Data: make([][]string, mat.Rows),
	}

	for i := 0; i < mat.Rows; i++ {
		jsonData.Data[i] = make([]string, mat.Cols)
		for j := 0; j < mat.Cols; j++ {
			jsonData.Data[i][j] = fmt.Sprintf("%v", mat.Data[i][j])
		}
	}

	return json.Marshal(jsonData)
}

// SaveToJSON сохраняет матрицу в JSON файл
func (mat *Matrix[T]) SaveToJSON(filename string) error {
	data, err := mat.MarshalJSON()
	if err != nil {
		return fmt.Errorf("ошибка маршалинга: %w", err)
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("ошибка записи файла: %w", err)
	}

	return nil
}

// Determinant вычисляет определитель матрицы методом Гаусса
func (m *Matrix[T]) Determinant() T {
	if m.Rows != m.Cols {
		return m.Data[0][0].Zero()
	}

	// Клонируем матрицу, чтобы не изменять исходную
	mat := m.Clone()
	n := mat.Rows
	det := mat.Data[0][0].One() // Начинаем с единицы

	for i := 0; i < n; i++ {
		if mat.Data[i][i].Equal(mat.Data[i][i].Zero()) {
			// Ищем ненулевой элемент в том же столбце ниже
			swapped := false
			for j := i + 1; j < n; j++ {
				if !mat.Data[j][i].Equal(mat.Data[j][i].Zero()) {
					// Меняем строки местами
					mat.Data[i], mat.Data[j] = mat.Data[j], mat.Data[i]
					det = det.Neg() // При перестановке строк меняем знак определителя
					swapped = true
					break
				}
			}
			if !swapped {
				return mat.Data[0][0].Zero() // Матрица вырождена
			}
		}

		// Умножаем определитель на диагональный элемент
		det = det.Mul(mat.Data[i][i])

		// Обнуляем элементы под диагональю
		for j := i + 1; j < n; j++ {
			if !mat.Data[j][i].Equal(mat.Data[j][i].Zero()) {
				factor, _ := mat.Data[j][i].Div(mat.Data[i][i])
				for k := i; k < n; k++ {
					mat.Data[j][k] = mat.Data[j][k].Sub(factor.Mul(mat.Data[i][k]))
				}
			}
		}
	}

	return det
}

// Rank вычисляет ранг матрицы методом Гаусса
func (m *Matrix[T]) Rank() int {
	// Клонируем матрицу, чтобы не изменять исходную
	mat := m.Clone()
	rank := 0
	rowCount := mat.Rows
	colCount := mat.Cols
	h := 0 // индекс текущей строки
	k := 0 // индекс текущего столбца

	for h < rowCount && k < colCount {
		// Находим строку с ненулевым элементом в столбце k
		i_max := h
		for i := h + 1; i < rowCount; i++ {
			if !mat.Data[i][k].Equal(mat.Data[i][k].Zero()) {
				i_max = i
				break
			}
		}

		if mat.Data[i_max][k].Equal(mat.Data[i_max][k].Zero()) {
			// Нет ненулевых элементов в этом столбце
			k++
			continue
		}

		// Меняем строки местами
		if i_max != h {
			mat.Data[h], mat.Data[i_max] = mat.Data[i_max], mat.Data[h]
		}

		// Обнуляем элементы в столбце k ниже строки h
		for i := h + 1; i < rowCount; i++ {
			if !mat.Data[i][k].Equal(mat.Data[i][k].Zero()) {
				factor, _ := mat.Data[i][k].Div(mat.Data[h][k])
				for j := k; j < colCount; j++ {
					mat.Data[i][j] = mat.Data[i][j].Sub(factor.Mul(mat.Data[h][j]))
				}
			}
		}

		rank++
		h++
		k++
	}

	return rank
}

// Inverse вычисляет обратную матрицу методом Гаусса-Жордана
func (m *Matrix[T]) Inverse() (*Matrix[T], error) {
	if m.Rows != m.Cols {
		return nil, errors.New("матрица должна быть квадратной")
	}

	n := m.Rows
	// Создаем расширенную матрицу [A|E]
	var zero T
	var one T = zero.One()
	augmented := NewMatrix[T](n, 2*n, zero)

	// Копируем исходную матрицу в левую часть
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			augmented.Data[i][j] = m.Data[i][j]
		}
		// Добавляем единичную матрицу справа
		augmented.Data[i][i+n] = one
	}

	// Прямой ход метода Гаусса
	for i := 0; i < n; i++ {
		pivot := augmented.Data[i][i]
		if pivot.Equal(zero) {
			// Ищем ненулевой элемент в том же столбце ниже
			swapped := false
			for j := i + 1; j < n; j++ {
				if !augmented.Data[j][i].Equal(zero) {
					// Меняем строки местами
					augmented.Data[i], augmented.Data[j] = augmented.Data[j], augmented.Data[i]
					pivot = augmented.Data[i][i]
					swapped = true
					break
				}
			}
			if !swapped {
				return nil, errors.New("матрица вырождена")
			}
		}

		// Делим строку на ведущий элемент
		for j := 0; j < 2*n; j++ {
			d, _ := augmented.Data[i][j].Div(pivot)
			augmented.Data[i][j] = d
		}

		// Обнуляем элементы в столбце i
		for k := 0; k < n; k++ {
			if k != i {
				factor := augmented.Data[k][i]
				for j := 0; j < 2*n; j++ {
					augmented.Data[k][j] = augmented.Data[k][j].Sub(factor.Mul(augmented.Data[i][j]))
				}
			}
		}
	}

	inverse := NewMatrix[T](n, n, zero)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			inverse.Data[i][j] = augmented.Data[i][j+n]
		}
	}

	return inverse, nil
}
