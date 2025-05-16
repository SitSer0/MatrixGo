package field

import (
	"fmt"
	"math/big"
)

// GF представляет элемент конечного поля GF(p)
type GF struct {
	value *big.Int
	p     *big.Int // характеристика поля (простое число)
}

// NewGF создает новый элемент конечного поля GF(p)
func NewGF(value int64, p int64) (GF, error) {
	if !isPrime(p) {
		return GF{}, fmt.Errorf("характеристика поля должна быть простым числом")
	}

	bigP := big.NewInt(p)
	val := new(big.Int).Mod(big.NewInt(value), bigP)

	return GF{
		value: val,
		p:     bigP,
	}, nil
}

// isPrime проверяет, является ли число простым
func isPrime(n int64) bool {
	if n <= 1 {
		return false
	}
	if n <= 3 {
		return true
	}
	if n%2 == 0 || n%3 == 0 {
		return false
	}

	for i := int64(5); i*i <= n; i += 6 {
		if n%i == 0 || n%(i+2) == 0 {
			return false
		}
	}
	return true
}

func (g GF) Add(other GF) GF {
	if g.p.Cmp(other.p) != 0 {
		panic("операции возможны только над элементами одного поля")
	}

	result := new(big.Int).Add(g.value, other.value)
	result.Mod(result, g.p)

	return GF{value: result, p: g.p}
}

func (g GF) Sub(other GF) GF {
	if g.p.Cmp(other.p) != 0 {
		panic("операции возможны только над элементами одного поля")
	}

	result := new(big.Int).Sub(g.value, other.value)
	result.Mod(result, g.p)

	return GF{value: result, p: g.p}
}

func (g GF) Mul(other GF) GF {
	if g.p.Cmp(other.p) != 0 {
		panic("операции возможны только над элементами одного поля")
	}

	result := new(big.Int).Mul(g.value, other.value)
	result.Mod(result, g.p)

	return GF{value: result, p: g.p}
}

func (g GF) Div(other GF) (GF, error) {
	if g.p.Cmp(other.p) != 0 {
		return GF{}, fmt.Errorf("операции возможны только над элементами одного поля")
	}

	if other.value.Sign() == 0 {
		return GF{}, fmt.Errorf("деление на ноль")
	}

	// Находим мультипликативно обратный элемент с помощью расширенного алгоритма Евклида
	inverse := new(big.Int)
	if inverse.ModInverse(other.value, g.p) == nil {
		return GF{}, fmt.Errorf("не существует мультипликативно обратного элемента")
	}

	result := new(big.Int).Mul(g.value, inverse)
	result.Mod(result, g.p)

	return GF{value: result, p: g.p}, nil
}

func (g GF) Neg() GF {
	result := new(big.Int).Neg(g.value)
	result.Mod(result, g.p)
	return GF{value: result, p: g.p}
}

func (g GF) Zero() GF {
	return GF{value: big.NewInt(0), p: g.p}
}

func (g GF) One() GF {
	return GF{value: big.NewInt(1), p: g.p}
}

func (g GF) Equal(other GF) bool {
	return g.p.Cmp(other.p) == 0 && g.value.Cmp(other.value) == 0
}

func (g GF) String() string {
	return fmt.Sprintf("%d (mod %d)", g.value, g.p)
}

// Проверка реализации интерфейса Field
var _ Field[GF] = GF{}
