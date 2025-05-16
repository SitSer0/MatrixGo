package field

import (
	"fmt"
	"math/big"
)

// Rational представляет рациональное число как отношение двух больших целых чисел
type Rational struct {
	num *big.Int // числитель
	den *big.Int // знаменатель
}

// NewRational создает новое рациональное число
func NewRational(num, den int64) Rational {
	r := Rational{
		num: big.NewInt(num),
		den: big.NewInt(den),
	}
	return r.normalize()
}

// NewRationalFromBig создает новое рациональное число из больших целых чисел
func NewRationalFromBig(num, den *big.Int) Rational {
	r := Rational{
		num: new(big.Int).Set(num),
		den: new(big.Int).Set(den),
	}
	return r.normalize()
}

// normalize приводит дробь к несократимому виду и обеспечивает положительный знаменатель
func (r Rational) normalize() Rational {
	if r.den.Sign() == 0 {
		panic("знаменатель не может быть равен нулю")
	}

	// Если знаменатель отрицательный, меняем знаки числителя и знаменателя
	if r.den.Sign() < 0 {
		r.num.Neg(r.num)
		r.den.Neg(r.den)
	}

	// Находим НОД и сокращаем дробь
	gcd := new(big.Int).GCD(nil, nil, r.num, r.den)
	if gcd.Cmp(big.NewInt(1)) > 0 {
		r.num.Div(r.num, gcd)
		r.den.Div(r.den, gcd)
	}

	return r
}

func (r Rational) Add(other Rational) Rational {
	// a/b + c/d = (ad + bc)/bd
	num := new(big.Int).Mul(r.num, other.den)
	tmp := new(big.Int).Mul(other.num, r.den)
	num.Add(num, tmp)
	den := new(big.Int).Mul(r.den, other.den)

	return NewRationalFromBig(num, den)
}

func (r Rational) Sub(other Rational) Rational {
	// a/b - c/d = (ad - bc)/bd
	num := new(big.Int).Mul(r.num, other.den)
	tmp := new(big.Int).Mul(other.num, r.den)
	num.Sub(num, tmp)
	den := new(big.Int).Mul(r.den, other.den)

	return NewRationalFromBig(num, den)
}

func (r Rational) Mul(other Rational) Rational {
	// (a/b) * (c/d) = (ac)/(bd)
	num := new(big.Int).Mul(r.num, other.num)
	den := new(big.Int).Mul(r.den, other.den)

	return NewRationalFromBig(num, den)
}

func (r Rational) Div(other Rational) (Rational, error) {
	if other.num.Sign() == 0 {
		return Rational{}, fmt.Errorf("деление на ноль")
	}
	// (a/b) / (c/d) = (ad)/(bc)
	num := new(big.Int).Mul(r.num, other.den)
	den := new(big.Int).Mul(r.den, other.num)

	return NewRationalFromBig(num, den), nil
}

func (r Rational) Neg() Rational {
	return NewRationalFromBig(new(big.Int).Neg(r.num), new(big.Int).Set(r.den))
}

func (r Rational) Zero() Rational {
	return NewRational(0, 1)
}

func (r Rational) One() Rational {
	return NewRational(1, 1)
}

func (r Rational) Equal(other Rational) bool {
	// Предполагаем, что обе дроби нормализованы
	return r.num.Cmp(other.num) == 0 && r.den.Cmp(other.den) == 0
}

func (r Rational) String() string {
	if r.den.Cmp(big.NewInt(1)) == 0 {
		return r.num.String()
	}
	return fmt.Sprintf("%s/%s", r.num.String(), r.den.String())
}

// Проверка реализации интерфейса Field
var _ Field[Rational] = Rational{}
