package entity_test

import (
	"math"
	"testing"

	"github.com/biraneves/fc-labs-weather/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTemperatureKelvin(t *testing.T) {
	tests := []struct {
		name          string
		input         float64
		expectedError string
	}{
		{
			name:  "valid kelvin temperature",
			input: 82.4,
		},
		{
			name:  "zero kelvin",
			input: 0.0,
		},
		{
			name:          "below absolute zero",
			input:         -10.0,
			expectedError: "temperature kelvin: below absolute zero",
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
			got, err := entity.NewTemperatureKelvin(tt.input)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Equal(t, entity.TemperatureKelvin{}, got)
			} else {
				require.NoError(t, err)
				assert.NotEqual(t, entity.TemperatureKelvin{}, got)
			}
		})
	}
}

func TestTemperatureKelvin_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		input float64
		want  bool
	}{
		{
			name:  "valid kelvin temperature",
			input: 540.2,
			want:  true,
		},
		{
			name:  "zero kelvin",
			input: 0.0,
			want:  true,
		},
		{
			name:  "below absolute zero",
			input: -2.3,
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
			v, _ := entity.NewTemperatureKelvin(tt.input)
			got := v.IsValid()
			require.Equal(t, tt.want, got)
		})
	}
}

func TestTemperatureKelvin_Value(t *testing.T) {
	tests := []struct {
		name  string
		input float64
		want  float64
		valid bool
	}{
		{
			name:  "valid kelvin temperature",
			input: 87.3,
			want:  87.3,
			valid: true,
		},
		{
			name:  "zero kelvin",
			input: 0,
			want:  0.0,
			valid: true,
		},
		{
			name:  "below absolute zero",
			input: -10.0,
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
			v, _ := entity.NewTemperatureKelvin(tt.input)
			got := v.Value()
			valid := v.IsValid()

			require.Equal(t, tt.want, got)
			require.Equal(t, tt.valid, valid)
		})
	}
}

func TestTemperatureKelvin_String(t *testing.T) {
	tests := []struct {
		name  string
		input float64
		want  string
	}{
		{
			name:  "valid kelvin temperature",
			input: 340.0,
			want:  "340,0 K",
		},
		{
			name:  "zero kelvin",
			input: 0.0,
			want:  "0,0 K",
		},
		{
			name:  "below absolute zero",
			input: -280.0,
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
			v, _ := entity.NewTemperatureKelvin(tt.input)
			got := v.String()
			require.Equal(t, tt.want, got)
		})
	}
}

func TestTemperatureKelvin_Equal(t *testing.T) {
	tests := []struct {
		name string
		a    float64
		b    float64
		want bool
	}{
		{
			name: "equal valid temperatures",
			a:    98.7,
			b:    98.7,
			want: true,
		},
		{
			name: "different valid temperatures",
			a:    78.3,
			b:    103.4,
			want: false,
		},
		{
			name: "first temperature invalid",
			a:    -10.0,
			b:    24.5,
			want: false,
		},
		{
			name: "second temperature invalid",
			a:    980.7,
			b:    -45.0,
			want: false,
		},
		{
			name: "two invalid temperatures",
			a:    -34.0,
			b:    math.Inf(1),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			va, _ := entity.NewTemperatureKelvin(tt.a)
			vb, _ := entity.NewTemperatureKelvin(tt.b)

			gotA := va.Equal(vb)
			gotB := vb.Equal(va)

			require.Equal(t, tt.want, gotA)
			require.Equal(t, tt.want, gotB)
		})
	}
}

func TestTemperatureKelvin_ToCelsius(t *testing.T) {
	tests := []struct {
		name  string
		input float64
		want  float64
	}{
		{
			name:  "valid kelvin temperature",
			input: 289.4,
			want:  16.25,
		},
		{
			name:  "water freezing temperature",
			input: 273.15,
			want:  0.0,
		},
		{
			name:  "water boiling temperature",
			input: 373.15,
			want:  100.0,
		},
		{
			name:  "absolute zero",
			input: 0.0,
			want:  -273.15,
		},
		{
			name:  "invalid temperature",
			input: -18.0,
			want:  math.NaN(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, _ := entity.NewTemperatureKelvin(tt.input)
			got := v.ToCelsius()
			if !v.IsValid() {
				require.True(t, math.IsNaN(got))
			} else {
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func TestTemperatureKelvin_ToFahrenheit(t *testing.T) {
	tests := []struct {
		name  string
		input float64
		want  float64
	}{
		{
			name:  "valid kelvin temperature",
			input: 328.4,
			want:  131.4,
		},
		{
			name:  "water freezing temperature",
			input: 273.15,
			want:  32.0,
		},
		{
			name:  "water boiling temperature",
			input: 373.15,
			want:  212.0,
		},
		{
			name:  "absolute zero",
			input: 0.0,
			want:  -459.67,
		},
		{
			name:  "invalid temperature",
			input: -10.0,
			want:  math.NaN(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, _ := entity.NewTemperatureKelvin(tt.input)
			got := v.ToFahrenheit()

			if !v.IsValid() {
				require.True(t, math.IsNaN(got))
			} else {
				require.True(t, entity.AlmostEqual(tt.want, got, EPSILON))
			}
		})
	}
}

func TestTemperatureKelvin_MarshalJSON(t *testing.T) {
	tests := []struct {
		name  string
		input float64
		want  []byte
	}{
		{
			name:  "valid temperature",
			input: 280.4,
			want:  []byte("280.4"),
		},
		{
			name:  "invalid temperature",
			input: -48.0,
			want:  []byte("null"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, _ := entity.NewTemperatureKelvin(tt.input)
			got, err := v.MarshalJSON()

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTemperatureKelvin_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name          string
		input         []byte
		want          float64
		valid         bool
		expectedError string
	}{
		{
			name:  "valid temperature",
			input: []byte(`18.9`),
			want:  18.9,
			valid: true,
		},
		{
			name:          "temperature as string",
			input:         []byte(`"240.0"`),
			want:          0.0,
			valid:         false,
			expectedError: "json: cannot unmarshal string",
		},
		{
			name:          "invalid temperature",
			input:         []byte(`-80.0`),
			want:          0.0,
			valid:         false,
			expectedError: "temperature kelvin: below absolute zero",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var v entity.TemperatureKelvin
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
