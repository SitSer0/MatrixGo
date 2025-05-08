package field

import (
	"fmt"
	"math"
	"strconv"
)

type Float64 float64

func (a Float64) Add(b Float64) Float64 { return a + b }
func (a Float64) Sub(b Float64) Float64 { return a - b }
func (a Float64) Mul(b Float64) Float64 { return a * b }
func (a Float64) Div(b Float64) (Float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("деление на ноль")
	}
	return a / b, nil
}
func (a Float64) Neg() Float64 { return -a }

func (a Float64) Zero() Float64 { return 0 }
func (a Float64) One() Float64  { return 1 }

func (f Float64) String() string {
	return fmt.Sprintf("%g", float64(f))
}

var _ Field[Float64] = Float64(0)

func (a Float64) Equal(b Float64) bool {
	const eps = 1e-9
	return math.Abs(float64(a-b)) < eps
}

func ParseFloat64(s string) (Float64, error) {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("не удалось преобразовать %q в float64", s)
	}
	return Float64(v), nil
}
