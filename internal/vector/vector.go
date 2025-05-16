package vector

import (
	"MatrixGo/internal/field"
	"errors"
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

func (v *Vector[T]) Add(v1 *Vector[T]) (*Vector[T], error) {
	if v.Len() != v1.Len() {
		return nil, errors.New("разные размеры векторов")
	}

	res := NewVector(v.Data)
	for i := 0; i < v.Len(); i++ {
		res.Data[i] = res.Data[i].Add(v1.Data[i])
	}
	return res, nil
}
