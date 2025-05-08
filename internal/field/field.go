package field

type Field[T any] interface {
	Add(other T) T
	Sub(other T) T
	Mul(other T) T
	Div(other T) (T, error) // деление может возвращать ошибку (деление на ноль)
	Neg() T                 // унарный минус
	Zero() T                // нулевой элемент поля
	One() T                 // единичный элемент поля
	Equal(other T) bool
}
