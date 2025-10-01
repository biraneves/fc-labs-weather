package entity_test

import (
	"math"
	"testing"

	"github.com/biraneves/fc-labs-weather/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTemperatureFahrenheit(t *testing.T) {
	tests := []struct {
		name          string
		input         float64
		expectedError string
	}{
		{
			name:  "valid fahrenheit temperature",
			input: 48.0,
		},
		{
			name:  "zero degree fahrenheit",
			input: 0.0,
		},
		{
			name:          "below absolute zero",
			input:         -500,
			expectedError: "temperature fahrenheit: below absolute zero",
		},
		{
			name:          "not a number",
			input:         math.NaN(),
			expectedError: "temperature: value is NaN",
		},
		{
			name:          "positive infinite",
			input:         math.Inf(1),
			expectedError: "temperature: value is infinite",
		},
		{
			name:          "negative infinite",
			input:         math.Inf(-1),
			expectedError: "temperature: value is infinite",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := entity.NewTemperatureFahrenheit(tt.input)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Equal(t, entity.TemperatureFahrenheit{}, got)
			} else {
				require.NoError(t, err)
				assert.NotEqual(t, entity.TemperatureFahrenheit{}, got)
			}
		})
	}
}

func TestTemperatureFahrenheit_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		input float64
		want  bool
	}{
		{
			name:  "valid fahrenheit temperature",
			input: 118.5,
			want:  true,
		},
		{
			name:  "zero degree fahrenheit",
			input: 0.0,
			want:  true,
		},
		{
			name:  "below absolute zero",
			input: -650.3,
			want:  false,
		},
		{
			name:  "not a number",
			input: math.NaN(),
			want:  false,
		},
		{
			name:  "positive infinite",
			input: math.Inf(1),
			want:  false,
		},
		{
			name:  "negative infinite",
			input: math.Inf(-1),
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, _ := entity.NewTemperatureFahrenheit(tt.input)
			got := v.IsValid()
			require.Equal(t, tt.want, got)
		})
	}
}

func TestTemperatureFahrenheit_Value(t *testing.T) {
	tests := []struct {
		name  string
		input float64
		want  float64
		valid bool
	}{
		{
			name:  "valid fahrenheit temperature",
			input: 98.3,
			want:  98.3,
			valid: true,
		},
		{
			name:  "zero degree fahrenheit",
			input: 0.0,
			want:  0.0,
			valid: true,
		},
		{
			name:  "below absolute zero",
			input: -820,
			want:  0.0,
			valid: false,
		},
		{
			name:  "not a number",
			input: math.NaN(),
			want:  0.0,
			valid: false,
		},
		{
			name:  "positive infinite",
			input: math.Inf(1),
			want:  0.0,
			valid: false,
		},
		{
			name:  "negative infinite",
			input: math.Inf(-1),
			want:  0.0,
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, _ := entity.NewTemperatureFahrenheit(tt.input)
			got := v.Value()
			valid := v.IsValid()

			require.Equal(t, tt.want, got)
			require.Equal(t, tt.valid, valid)
		})
	}
}

func TestTemperatureFahrenheit_String(t *testing.T) {
	tests := []struct {
		name  string
		input float64
		want  string
	}{
		{
			name:  "valid fahrenheit temperature",
			input: 357.46,
			want:  "357,5 °F",
		},
		{
			name:  "zero degree fahrenheit",
			input: 0.0,
			want:  "0,0 °F",
		},
		{
			name:  "below absolute zero",
			input: -684.8,
			want:  "",
		},
		{
			name:  "not a number",
			input: math.NaN(),
			want:  "",
		},
		{
			name:  "positive infinite",
			input: math.Inf(1),
			want:  "",
		},
		{
			name:  "negative infinite",
			input: math.Inf(-1),
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, _ := entity.NewTemperatureFahrenheit(tt.input)
			got := v.String()
			require.Equal(t, tt.want, got)
		})
	}
}

func TestTemperatureFahrenheit_Equal(t *testing.T) {
	tests := []struct {
		name string
		a    float64
		b    float64
		want bool
	}{
		{
			name: "equal valid temperatures",
			a:    18.3,
			b:    18.3,
			want: true,
		},
		{
			name: "different valid temperatures",
			a:    89.4,
			b:    -34.7,
			want: false,
		},
		{
			name: "first temperature invalid",
			a:    -2380.4,
			b:    54.3,
			want: false,
		},
		{
			name: "second temperature invalid",
			a:    98.9,
			b:    math.NaN(),
			want: false,
		},
		{
			name: "two invalid temperatures",
			a:    math.Inf(1),
			b:    -500.0,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			va, _ := entity.NewTemperatureFahrenheit(tt.a)
			vb, _ := entity.NewTemperatureFahrenheit(tt.b)

			gotA := va.Equal(vb)
			gotB := vb.Equal(va)

			require.Equal(t, tt.want, gotA)
			require.Equal(t, tt.want, gotB)
		})
	}
}

func TestTemperatureFahrenheit_ToCelsius(t *testing.T) {
	tests := []struct {
		name  string
		input float64
		want  float64
	}{
		{
			name:  "valid fahrenheit temperature",
			input: 18.4,
			want:  -7.55555,
		},
		{
			name:  "water freezing temperature",
			input: 32.0,
			want:  0.0,
		},
		{
			name:  "water boiling temperature",
			input: 212.0,
			want:  100.0,
		},
		{
			name:  "absolute zero",
			input: -459.67,
			want:  -273.15,
		},
		{
			name:  "invalid temperature",
			input: -460.0,
			want:  math.NaN(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, _ := entity.NewTemperatureFahrenheit(tt.input)
			got := v.ToCelsius()

			if !v.IsValid() {
				require.True(t, math.IsNaN(got))
			} else {
				require.True(t, entity.AlmostEqual(tt.want, got, EPSILON))
			}
		})
	}
}

func TestTemperatureFahrenheit_ToKelvin(t *testing.T) {
	tests := []struct {
		name  string
		input float64
		want  float64
	}{
		{
			name:  "valid fahrenheit temperature",
			input: 48.6,
			want:  282.37,
		},
		{
			name:  "water freezing temperature",
			input: 32.0,
			want:  273.15,
		},
		{
			name:  "water boiling temperature",
			input: 212.0,
			want:  373.15,
		},
		{
			name:  "absolute zero",
			input: -459.67,
			want:  0.0,
		},
		{
			name:  "invalid temperature",
			input: -502.8,
			want:  math.NaN(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, _ := entity.NewTemperatureFahrenheit(tt.input)
			got := v.ToKelvin()

			if !v.IsValid() {
				require.True(t, math.IsNaN(got))
			} else {
				require.True(t, entity.AlmostEqual(tt.want, got, EPSILON))
			}
		})
	}
}

func TestTemperatureFahrenheit_MarshalJSON(t *testing.T) {
	tests := []struct {
		name  string
		input float64
		want  []byte
	}{
		{
			name:  "valid fahrenheit temperature",
			input: 98.4,
			want:  []byte("98.4"),
		},
		{
			name:  "invalid temperature",
			input: -840,
			want:  []byte("null"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, _ := entity.NewTemperatureFahrenheit(tt.input)
			got, err := v.MarshalJSON()
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTemperatureFahrenheit_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name          string
		input         []byte
		want          float64
		valid         bool
		expectedError string
	}{
		{
			name:  "valid fahrenheit temperature",
			input: []byte(`98.6`),
			want:  98.6,
			valid: true,
		},
		{
			name:          "temperature as string",
			input:         []byte(`"45.6"`),
			want:          0.0,
			valid:         false,
			expectedError: "json: cannot unmarshal string",
		},
		{
			name:          "invalid temperature",
			input:         []byte(`-840`),
			want:          0.0,
			valid:         false,
			expectedError: "temperature fahrenheit: below absolute zero",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var v entity.TemperatureFahrenheit
			err := v.UnmarshalJSON(tt.input)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.False(t, v.IsValid())
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, v.Value())
				assert.Equal(t, tt.valid, v.IsValid())
			}
		})
	}
}
