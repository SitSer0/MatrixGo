package io

import (
	"MatrixGo/internal/field"
	"MatrixGo/internal/matrix"
	"encoding/csv"
	"fmt"
	"os"
)

type CSVParser[T any] func(string) (T, error)

func ReadMatrixFromCSV[T field.Field[T]](filename string, parse CSVParser[T]) (*matrix.Matrix[T], error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть файл: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	rows, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения CSV: %w", err)
	}

	if len(rows) == 0 || len(rows[0]) == 0 {
		return nil, fmt.Errorf("пустой CSV")
	}

	rowCount := len(rows)
	colCount := len(rows[0])

	var zero T
	mat := matrix.NewMatrix[T](rowCount, colCount, zero)

	for i := 0; i < rowCount; i++ {
		if len(rows[i]) != colCount {
			return nil, fmt.Errorf("разная длина строк в CSV (строка %d)", i)
		}
		for j := 0; j < colCount; j++ {
			val, err := parse(rows[i][j])
			if err != nil {
				return nil, fmt.Errorf("ошибка парсинга [%d][%d]: %w", i, j, err)
			}
			mat.Data[i][j] = val
		}
	}

	return mat, nil
}

func WriteMatrixToCSV[T field.Field[T]](filename string, mat *matrix.Matrix[T]) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("не удалось создать файл: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for i := 0; i < mat.Rows; i++ {
		row := make([]string, mat.Cols)
		for j := 0; j < mat.Cols; j++ {
			row[j] = fmt.Sprintf("%v", mat.Data[i][j]) // можно переопределить String() у T
		}
		err := writer.Write(row)
		if err != nil {
			return fmt.Errorf("ошибка записи строки: %w", err)
		}
	}

	return nil
}
