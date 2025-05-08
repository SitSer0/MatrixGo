package main

import (
	"MatrixGo/internal/field"
	"MatrixGo/internal/matrix"
	"fmt"
)

func main() {
	a := field.Float64(3.14)
	b := field.Float64(2.0)
	fmt.Println("Сумма:", a.Add(b))
	m := matrix.NewMatrix[field.Complex](3, 3, field.Complex{Re: 1})
	n := matrix.NewMatrix[field.Complex](3, 3, field.Complex{Re: 1})
	c, err := m.Mul(n)
	fmt.Println(c.Data, err)
}
