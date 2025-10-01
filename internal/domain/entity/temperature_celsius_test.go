package entity_test

import (
	"math"
	"testing"

	"github.com/biraneves/fc-labs-weather/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const EPSILON = 0.1

func TestNewTemperatureCelsius(t *testing.T) {
	tests := []struct {
		name          string
		input         float64
		expectedError string
	}{
		{
			name:  "valid celsius temperature",
			input: 25.0,
		},
		{
			name:  "zero degree celsius",
			input: 0.0,
		},
		{
			name:          "below absolute zero",
			input:         -280.0,
			expectedError: "temperature celsius: below absolute zero",
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
			got, err := entity.NewTemperatureCelsius(tt.input)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Equal(t, entity.TemperatureCelsius{}, got)
			} else {
				require.NoError(t, err)
				assert.NotEqual(t, entity.TemperatureCelsius{}, got)
			}
		})
	}
}

func TestTemperatureCelsius_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		input float64
		want  bool
	}{
		{
			name:  "valid celsius temperature",
			input: 25.0,
			want:  true,
		},
		{
			name:  "zero degree celsius",
			input: 0.0,
			want:  true,
		},
		{
			name:  "below absolute zero",
			input: -300.0,
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
			v, _ := entity.NewTemperatureCelsius(tt.input)
			got := v.IsValid()
			require.Equal(t, tt.want, got)
		})
	}
}

func TestTemperatureCelsius_Value(t *testing.T) {
	tests := []struct {
		name  string
		input float64
		want  float64
		valid bool
	}{
		{
			name:  "valid celsius temperature",
			input: 18.4,
			want:  18.4,
			valid: true,
		},
		{
			name:  "zero degree celsius",
			input: 0.0,
			want:  0.0,
			valid: true,
		},
		{
			name:  "below absolute zero",
			input: -280.3,
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
			v, _ := entity.NewTemperatureCelsius(tt.input)
			got := v.Value()
			valid := v.IsValid()

			require.Equal(t, tt.want, got)
			require.Equal(t, tt.valid, valid)
		})
	}
}

func TestTemperatureCelsius_String(t *testing.T) {
	tests := []struct {
		name  string
		input float64
		want  string
	}{
		{
			name:  "valid celsius temperature",
			input: 25.0,
			want:  "25,0 °C",
		},
		{
			name:  "zero degree celsius",
			input: 0.0,
			want:  "0,0 °C",
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
			v, _ := entity.NewTemperatureCelsius(tt.input)
			got := v.String()
			require.Equal(t, tt.want, got)
		})
	}
}

func TestTemperatureCelsius_Equal(t *testing.T) {
	tests := []struct {
		name string
		a    float64
		b    float64
		want bool
	}{
		{
			name: "equal valid temperatures",
			a:    25.0,
			b:    25,
			want: true,
		},
		{
			name: "different valid temperatures",
			a:    18.4,
			b:    -24.3,
			want: false,
		},
		{
			name: "first temperature invalid",
			a:    -280,
			b:    340.5,
			want: false,
		},
		{
			name: "second temperature invalid",
			a:    98.3,
			b:    math.NaN(),
			want: false,
		},
		{
			name: "two invalid temperatures",
			a:    -340.0,
			b:    math.Inf(1),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			va, _ := entity.NewTemperatureCelsius(tt.a)
			vb, _ := entity.NewTemperatureCelsius(tt.b)

			gotA := va.Equal(vb)
			gotB := vb.Equal(va)

			require.Equal(t, tt.want, gotA)
			require.Equal(t, tt.want, gotB)
		})
	}
}

func TestTemperatureCelsius_ToFahrenheit(t *testing.T) {
	tests := []struct {
		name  string
		input float64
		want  float64
	}{
		{
			name:  "valid celsius temperature",
			input: 25.0,
			want:  77.0,
		},
		{
			name:  "water freezing temperature",
			input: 0.0,
			want:  32.0,
		},
		{
			name:  "water boiling temperature",
			input: 100.0,
			want:  212.0,
		},
		{
			name:  "absolute zero",
			input: -273.15,
			want:  -459.67,
		},
		{
			name:  "invalid temperature",
			input: -300,
			want:  math.NaN(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, _ := entity.NewTemperatureCelsius(tt.input)
			got := v.ToFahrenheit()

			if !v.IsValid() {
				require.True(t, math.IsNaN(got))
			} else {
				require.True(t, entity.AlmostEqual(tt.want, got, EPSILON))
			}
		})
	}
}

func TestTemperatureCelsius_ToKelvin(t *testing.T) {
	tests := []struct {
		name  string
		input float64
		want  float64
	}{
		{
			name:  "valid temperature",
			input: 19.0,
			want:  292.15,
		},
		{
			name:  "water freezing temperature",
			input: 0,
			want:  273.15,
		},
		{
			name:  "water boiling temperature",
			input: 100,
			want:  373.15,
		},
		{
			name:  "absolute zero",
			input: -273.15,
			want:  0,
		},
		{
			name:  "invalid temperature",
			input: -450,
			want:  math.NaN(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, _ := entity.NewTemperatureCelsius(tt.input)
			got := v.ToKelvin()

			if !v.IsValid() {
				require.True(t, math.IsNaN(got))
			} else {
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func TestTemperatureCelsius_MarshalJSON(t *testing.T) {
	tests := []struct {
		name  string
		input float64
		want  []byte
	}{
		{
			name:  "valid celsius temperature",
			input: 28.3,
			want:  []byte("28.3"),
		},
		{
			name:  "invalid temperature",
			input: -540.8,
			want:  []byte("null"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, _ := entity.NewTemperatureCelsius(tt.input)
			got, err := v.MarshalJSON()

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTemperatureCelsius_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name          string
		input         []byte
		want          float64
		valid         bool
		expectedError string
	}{
		{
			name:  "valid celsius temperature",
			input: []byte(`18.5`),
			want:  18.5,
			valid: true,
		},
		{
			name:          "temperature as string",
			input:         []byte(`"23.4"`),
			want:          0.0,
			valid:         false,
			expectedError: "json: cannot unmarshal string",
		},
		{
			name:          "invalid temperature",
			input:         []byte(`-280`),
			want:          0.0,
			valid:         false,
			expectedError: "temperature celsius: below absolute zero",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var v entity.TemperatureCelsius
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
