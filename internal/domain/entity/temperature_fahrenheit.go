package entity

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
)

type TemperatureFahrenheit struct {
	value float64
	valid bool
}

const absoluteZeroFahrenheit = -459.67

var (
	ErrTemFahrenheitBelowAbsZero = fmt.Errorf("temperature fahrenheit: below absolute zero (%.2f °F)", absoluteZeroFahrenheit)
)

func NewTemperatureFahrenheit(tempFahrenheit float64) (TemperatureFahrenheit, error) {
	if math.IsNaN(tempFahrenheit) {
		return TemperatureFahrenheit{}, ErrTempNan
	}

	if math.IsInf(tempFahrenheit, 0) {
		return TemperatureFahrenheit{}, ErrTempInf
	}

	if tempFahrenheit < absoluteZeroFahrenheit {
		return TemperatureFahrenheit{}, ErrTemFahrenheitBelowAbsZero
	}

	return TemperatureFahrenheit{tempFahrenheit, true}, nil
}

func (t TemperatureFahrenheit) IsValid() bool {
	return t.valid
}

func (t TemperatureFahrenheit) Value() float64 {
	return t.value
}

func (t TemperatureFahrenheit) String() string {
	if !t.IsValid() {
		return ""
	}

	return strings.Replace(fmt.Sprintf("%.1f °F", t.value), ".", ",", 1)
}

func (t TemperatureFahrenheit) Equal(o TemperatureFahrenheit) bool {
	if !t.IsValid() && !o.IsValid() {
		return false
	}

	return t.value == o.value && t.valid == o.valid
}

func (t TemperatureFahrenheit) ToCelsius() float64 {
	if !t.IsValid() {
		return math.NaN()
	}

	c := 5.0 * (t.value - 32) / 9.0
	return math.Round(c*10) / 10
}

func (t TemperatureFahrenheit) ToKelvin() float64 {
	if !t.IsValid() {
		return math.NaN()
	}

	k := (5.0 * (t.value - 32.0) / 9.0) + 273.15
	return math.Round(k*10) / 10
}

func (t TemperatureFahrenheit) MarshalJSON() ([]byte, error) {
	if !t.IsValid() {
		return []byte("null"), nil
	}

	v := math.Round(t.value*10) / 10
	formatted := strconv.FormatFloat(v, 'f', 1, 64)
	return []byte(formatted), nil
}

func (t *TemperatureFahrenheit) UnmarshalJSON(b []byte) error {
	var f float64
	err := json.Unmarshal(b, &f)
	if err != nil {
		return err
	}

	v, err := NewTemperatureFahrenheit(f)
	if err != nil {
		return err
	}

	*t = v
	return nil
}
