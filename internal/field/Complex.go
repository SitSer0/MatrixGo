package field

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Complex struct {
	Re float64
	Im float64
}

func NewComplex(a, b float64) *Complex {
	return &Complex{Re: a, Im: b}
}

func (a Complex) Add(b Complex) Complex {
	return Complex{Re: a.Re + b.Re, Im: a.Im + b.Im}
}

func (a Complex) Sub(b Complex) Complex {
	return Complex{Re: a.Re - b.Re, Im: a.Im - b.Im}
}

func (a Complex) Mul(b Complex) Complex {
	return Complex{Re: a.Re*b.Re - a.Im*b.Im, Im: a.Re*b.Im + a.Im*b.Re}
}

func (a Complex) Div(b Complex) (Complex, error) {
	denominator := b.Re*b.Re + b.Im*b.Im
	if denominator == 0 {
		return Complex{}, errors.New("деление на ноль")
	}
	re := (a.Re*b.Re + a.Im*b.Im) / denominator
	im := (a.Im*b.Re - a.Re*b.Im) / denominator
	return Complex{Re: re, Im: im}, nil
}

func (a Complex) Neg() Complex {
	return Complex{Re: -a.Re, Im: -a.Im}
}

func (a Complex) Zero() Complex {
	return Complex{Re: 0, Im: 0}
}

func (a Complex) One() Complex {
	return Complex{Re: 1, Im: 0}
}

func (a Complex) Equal(b Complex) bool {
	const eps = 1e-9
	return math.Abs(a.Re-b.Re) < eps && math.Abs(a.Im-b.Im) < eps
}

func (c Complex) String() string {
	switch {
	case c.Re == 0 && c.Im == 0:
		return "0"
	case c.Re == 0:
		return fmt.Sprintf("%gi", c.Im)
	case c.Im == 0:
		return fmt.Sprintf("%g", c.Re)
	case c.Im > 0:
		return fmt.Sprintf("%g+%gi", c.Re, c.Im)
	default:
		return fmt.Sprintf("%g%gi", c.Re, c.Im) // знак уже в числе
	}
}

var _ Field[Complex] = Complex{0, 0}

// ParseComplex разбирает строки вида:
// "3+4i", "-2.1-5.3i", "5", "5i", "-i", "i"
func ParseComplex(s string) (Complex, error) {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, " ", "")

	// Обработка специальных случаев
	switch s {
	case "":
		return Complex{Re: 0, Im: 0}, nil
	case "i":
		return Complex{Re: 0, Im: 1}, nil
	case "-i":
		return Complex{Re: 0, Im: -1}, nil
	case "+i":
		return Complex{Re: 0, Im: 1}, nil
	}

	// Если строка заканчивается на i
	if strings.HasSuffix(s, "i") {
		// Убираем i в конце
		s = s[:len(s)-1]

		// Если после удаления i строка пустая, значит было просто i
		if s == "" {
			return Complex{Re: 0, Im: 1}, nil
		}

		// Если после удаления i остался только знак
		if s == "+" {
			return Complex{Re: 0, Im: 1}, nil
		}
		if s == "-" {
			return Complex{Re: 0, Im: -1}, nil
		}

		// Парсим как мнимую часть
		im, err := strconv.ParseFloat(s, 64)
		if err == nil {
			return Complex{Re: 0, Im: im}, nil
		}

		// Ищем + или - в середине строки (не в начале)
		plusIndex := -1
		minusIndex := -1
		for i := 1; i < len(s); i++ {
			if s[i] == '+' {
				plusIndex = i
				break
			}
			if s[i] == '-' {
				minusIndex = i
				break
			}
		}

		if plusIndex > 0 || minusIndex > 0 {
			splitIndex := plusIndex
			if minusIndex > 0 && (plusIndex == -1 || minusIndex < plusIndex) {
				splitIndex = minusIndex
			}

			// Парсим действительную часть
			re, err := strconv.ParseFloat(s[:splitIndex], 64)
			if err != nil {
				return Complex{}, fmt.Errorf("ошибка парсинга действительной части: %q", s[:splitIndex])
			}

			// Парсим мнимую часть
			im, err := strconv.ParseFloat(s[splitIndex:], 64)
			if err != nil {
				return Complex{}, fmt.Errorf("ошибка парсинга мнимой части: %q", s[splitIndex:])
			}

			return Complex{Re: re, Im: im}, nil
		}

		// Если нет разделителя, вся строка - коэффициент при i
		im, err = strconv.ParseFloat(s, 64)
		if err != nil {
			return Complex{}, fmt.Errorf("ошибка парсинга мнимой части: %q", s)
		}
		return Complex{Re: 0, Im: im}, nil
	}

	// Если нет i - это действительное число
	re, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return Complex{}, fmt.Errorf("ошибка парсинга действительного числа: %q", s)
	}
	return Complex{Re: re, Im: 0}, nil
}
