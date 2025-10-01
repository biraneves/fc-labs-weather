package entity

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
)

type TemperatureKelvin struct {
	value float64
	valid bool
}

const absoluteZeroKelvin = 0.0

var (
	ErrTempKelvinBelowAbsZero = fmt.Errorf("temperature kelvin: below absolute zero (%.1f K)", absoluteZeroKelvin)
)

func NewTemperatureKelvin(tempKelvin float64) (TemperatureKelvin, error) {
	if math.IsNaN(tempKelvin) {
		return TemperatureKelvin{}, ErrTempNan
	}

	if math.IsInf(tempKelvin, 0) {
		return TemperatureKelvin{}, ErrTempInf
	}

	if tempKelvin < absoluteZeroKelvin {
		return TemperatureKelvin{}, ErrTempKelvinBelowAbsZero
	}
	return TemperatureKelvin{tempKelvin, true}, nil
}

func (t TemperatureKelvin) IsValid() bool {
	return t.valid
}

func (t TemperatureKelvin) Value() float64 {
	return t.value
}

func (t TemperatureKelvin) String() string {
	if !t.IsValid() {
		return ""
	}

	return strings.Replace(fmt.Sprintf("%.1f K", t.value), ".", ",", 1)
}

func (t TemperatureKelvin) Equal(o TemperatureKelvin) bool {
	if !t.IsValid() && !o.IsValid() {
		return false
	}

	return t.value == o.value && t.valid == o.valid
}

func (t TemperatureKelvin) ToCelsius() float64 {
	if !t.IsValid() {
		return math.NaN()
	}

	return t.value - 273.15
}

func (t TemperatureKelvin) ToFahrenheit() float64 {
	if !t.IsValid() {
		return math.NaN()
	}

	f := (9.0 * (t.value - 273.15) / 5.0) + 32.0
	return math.Round(f*10) / 10
}

func (t TemperatureKelvin) MarshalJSON() ([]byte, error) {
	if !t.IsValid() {
		return []byte("null"), nil
	}

	v := math.Round(t.value*10) / 10
	formatted := strconv.FormatFloat(v, 'f', 1, 64)

	return []byte(formatted), nil
}

func (t *TemperatureKelvin) UnmarshalJSON(b []byte) error {
	var f float64
	err := json.Unmarshal(b, &f)
	if err != nil {
		return err
	}

	v, err := NewTemperatureKelvin(f)
	if err != nil {
		return err
	}

	*t = v
	return nil
}
