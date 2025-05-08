package vector

import (
	"MatrixGo/internal/field"
	"fmt"
)

type Vector[T field.Field[T]] struct {
	Size int
	Data []T
}

func NewVector[T field.Field[T]](data []T) *Vector[T] {
	return &Vector[T]{
		Size: len(data),
		Data: data,
	}
}

func (v *Vector[T]) Len() int {
	return v.Size
}

func (v *Vector[T]) At(index int) (T, error) {
	if index < 0 || index >= v.Size {
		return v.Data[0].Zero(), fmt.Errorf("индекс за пределами")
	}
	return v.Data[index], nil
}
