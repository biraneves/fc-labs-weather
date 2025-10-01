package entity

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

type TemperatureCelsius struct {
	value float64
	valid bool
}

const absoluteZeroCelsius = -273.15

var (
	ErrTempCelsiusBelowAbsZero = fmt.Errorf("temperature celsius: below absolute zero (%.2f °C)", absoluteZeroCelsius)
	ErrTempNan                 = errors.New("temperature: value is NaN")
	ErrTempInf                 = errors.New("temperature: value is infinite")
)

func NewTemperatureCelsius(tempCelsius float64) (TemperatureCelsius, error) {
	if math.IsNaN(tempCelsius) {
		return TemperatureCelsius{}, ErrTempNan
	}

	if math.IsInf(tempCelsius, 0) {
		return TemperatureCelsius{}, ErrTempInf
	}

	if tempCelsius < absoluteZeroCelsius {
		return TemperatureCelsius{}, ErrTempCelsiusBelowAbsZero
	}

	return TemperatureCelsius{tempCelsius, true}, nil
}

func (t TemperatureCelsius) IsValid() bool {
	return t.valid
}

func (t TemperatureCelsius) Value() float64 {
	return t.value
}

func (t TemperatureCelsius) String() string {
	if !t.IsValid() {
		return ""
	}

	return strings.Replace(fmt.Sprintf("%.1f °C", t.value), ".", ",", 1)
}

func (t TemperatureCelsius) Equal(o TemperatureCelsius) bool {
	if !t.IsValid() && !o.IsValid() {
		return false
	}

	return t.value == o.value
}

func (t TemperatureCelsius) ToFahrenheit() float64 {
	if !t.IsValid() {
		return math.NaN()
	}

	f := (9 * t.value / 5) + 32.0
	return math.Round(f*10) / 10
}

func (t TemperatureCelsius) ToKelvin() float64 {
	if !t.IsValid() {
		return math.NaN()
	}

	k := t.value - absoluteZeroCelsius
	return k
}

func (t TemperatureCelsius) MarshalJSON() ([]byte, error) {
	if !t.IsValid() {
		return []byte("null"), nil
	}

	v := math.Round(t.value*10) / 10
	formatted := strconv.FormatFloat(v, 'f', 1, 64)

	return []byte(formatted), nil
}

func (t *TemperatureCelsius) UnmarshalJSON(b []byte) error {
	var f float64
	err := json.Unmarshal(b, &f)
	if err != nil {
		return err
	}

	v, err := NewTemperatureCelsius(f)
	if err != nil {
		return err
	}

	*t = v
	return nil
}
