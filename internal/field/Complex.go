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

	switch s {
	case "i":
		return Complex{Re: 0, Im: 1}, nil
	case "-i":
		return Complex{Re: 0, Im: -1}, nil
	}

	if strings.HasSuffix(s, "i") {

		sWithoutI := s[:len(s)-1]

		lastSign := -1
		for i := len(sWithoutI) - 1; i > 0; i-- {
			if sWithoutI[i] == '+' || sWithoutI[i] == '-' {
				lastSign = i
				break
			}
		}

		if lastSign != -1 {
			rePartStr := sWithoutI[:lastSign]
			imPartStr := sWithoutI[lastSign:]

			rePart, err1 := strconv.ParseFloat(rePartStr, 64)
			imPart, err2 := strconv.ParseFloat(imPartStr, 64)
			if err1 != nil || err2 != nil {
				return Complex{}, fmt.Errorf("ошибка парсинга комплексного числа: %q", s)
			}
			return Complex{Re: rePart, Im: imPart}, nil
		}

		imPart, err := strconv.ParseFloat(sWithoutI, 64)
		if err != nil {
			return Complex{}, fmt.Errorf("ошибка парсинга мнимой части: %q", s)
		}
		return Complex{Re: 0, Im: imPart}, nil
	}

	rePart, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return Complex{}, fmt.Errorf("ошибка парсинга вещественной части: %q", s)
	}
	return Complex{Re: rePart, Im: 0}, nil
}
