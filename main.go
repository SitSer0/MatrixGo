package main

import (
	"MatrixGo/internal/field"
	"MatrixGo/internal/matrix"
	"fmt"
)

func main() {
	fmt.Println("1. Единичная матрица 3x3:")
	identity := matrix.Eye(3, field.Float64(0), field.Float64(1))
	fmt.Println(identity)
	fmt.Println()

	data := [][]field.Float64{
		{field.Float64(1), field.Float64(2), field.Float64(3)},
		{field.Float64(4), field.Float64(5), field.Float64(6)},
		{field.Float64(7), field.Float64(8), field.Float64(9)},
	}
	m, _ := matrix.FromSlice(data)
	fmt.Println("2. Матрица из слайса:")
	fmt.Println(m)
	fmt.Println()

	fmt.Println("3. Транспонированная матрица:")
	fmt.Println(m.Transpose())
	fmt.Println()

	c1 := matrix.NewMatrix[field.Complex](2, 2, field.Complex{Re: 1, Im: 1})
	c2 := matrix.NewMatrix[field.Complex](2, 2, field.Complex{Re: 0, Im: 2})
	fmt.Println("4. Матрица с комплексными числами:")
	fmt.Println(c1)
	fmt.Println()

	fmt.Println("5. Сложение комплексных матриц:")
	sum, _ := c1.Add(c2)
	fmt.Println(sum)

	fmt.Println("\n6. Сохранение матрицы в JSON:")
	err := m.SaveToJSON("matrix.json")
	if err != nil {
		fmt.Printf("Ошибка сохранения: %v\n", err)
	} else {
		fmt.Println("Матрица успешно сохранена в matrix.json")
	}

	fmt.Println("1. Матрица с рациональными числами:")
	ratData := [][]field.Rational{
		{field.NewRational(1, 2), field.NewRational(2, 3)},
		{field.NewRational(3, 4), field.NewRational(4, 5)},
	}
	ratMatrix, _ := matrix.FromSlice(ratData)
	fmt.Println(ratMatrix)
	fmt.Println()

	fmt.Println("2. Умножение матриц с рациональными числами:")
	ratMatrix2 := matrix.NewMatrix[field.Rational](2, 2, field.NewRational(1, 1))
	product, _ := ratMatrix.Mul(ratMatrix2)
	fmt.Println(product)
	fmt.Println()

	fmt.Println("3. Матрица над конечным полем GF(7):")
	gf7_1, _ := field.NewGF(3, 7)
	gf7_2, _ := field.NewGF(5, 7)
	gf7_3, _ := field.NewGF(2, 7)
	gf7_4, _ := field.NewGF(6, 7)

	gfData := [][]field.GF{
		{gf7_1, gf7_2},
		{gf7_3, gf7_4},
	}
	gfMatrix, _ := matrix.FromSlice(gfData)
	fmt.Println(gfMatrix)
	fmt.Println()

	fmt.Println("4. Сложение матриц в GF(7):")
	gfMatrix2 := matrix.NewMatrix[field.GF](2, 2, gf7_1)
	gfSum, err := gfMatrix.Add(gfMatrix2)
	fmt.Println(gfSum)
	fmt.Println()

	fmt.Println("5. Матрица с большими рациональными числами:")
	bigRat1 := field.NewRational(1234567890, 987654321)
	bigRatMatrix := matrix.NewMatrix[field.Rational](2, 2, bigRat1)
	fmt.Println(bigRatMatrix)
	fmt.Println()

	fmt.Println("6. Сохранение матрицы с рациональными числами в JSON:")
	err = ratMatrix.SaveToJSON("rational_matrix.json")
	if err != nil {
		fmt.Printf("Ошибка сохранения: %v\n", err)
	} else {
		fmt.Println("Матрица успешно сохранена в rational_matrix.json")
	}
}
