package entity_test

import (
	"testing"

	"github.com/biraneves/fc-labs-weather/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCep(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		want          string
		expectedError string
	}{
		{
			name:  "valid zip code",
			input: "07190050",
			want:  "07190050",
		},
		{
			name:  "leading and trailing white spaces",
			input: "   07190050  \t\n   \r",
			want:  "07190050",
		},
		{
			name:          "empty string",
			input:         "",
			want:          "",
			expectedError: "zip code: empty",
		},
		{
			name:          "less than eight digits",
			input:         "7190050",
			want:          "",
			expectedError: "zip code: invalid length",
		},
		{
			name:          "more than eight digits",
			input:         "071900500",
			want:          "",
			expectedError: "zip code: invalid length",
		},
		{
			name:          "invalid characters",
			input:         "7190-050",
			want:          "",
			expectedError: "zip code: invalid characters - only digits allowed",
		},
		{
			name:          "repeated digits",
			input:         "11111111",
			want:          "",
			expectedError: "zip code: eight equal digits not allowed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := entity.NewCep(tt.input)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Empty(t, got)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, string(got))
			}
		})
	}
}

func TestCep_String(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "valid zip code",
			input: "07190050",
			want:  "07190050",
		},
		{
			name:  "leading and trailing white spaces",
			input: "   07190050  ",
			want:  "07190050",
		},
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
		{
			name:  "less than eight digits",
			input: "7190050",
			want:  "",
		},
		{
			name:  "more than eight digits",
			input: "071900500",
			want:  "",
		},
		{
			name:  "invalid characters",
			input: "7190-050",
			want:  "",
		},
		{
			name:  "repeated digits",
			input: "11111111",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := entity.NewCep(tt.input)
			got := c.String()
			require.Equal(t, tt.want, got)
		})
	}
}

func TestCep_IsZero(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "empty zip code",
			input: "",
			want:  true,
		},
		{
			name:  "valid zip code",
			input: "12345678",
			want:  false,
		},
		{
			name:  "invalid zip code",
			input: "1234567",
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := entity.NewCep(tt.input)
			got := c.IsZero()
			require.Equal(t, tt.want, got)
		})
	}
}

func TestCep_Equal(t *testing.T) {
	tests := []struct {
		name string
		a    string
		b    string
		want bool
	}{
		{
			name: "equal valid zip codes",
			a:    "12345678",
			b:    "12345678",
			want: true,
		},
		{
			name: "different valid zip codes",
			a:    "12345678",
			b:    "23456789",
			want: false,
		},
		{
			name: "first zip code invalid",
			a:    "1234567",
			b:    "12345678",
			want: false,
		},
		{
			name: "second zip code invalid",
			a:    "12345678",
			b:    "1234567",
			want: false,
		},
		{
			name: "two invalid zip codes",
			a:    "1234567",
			b:    "11111111",
			want: false,
		},
		{
			name: "two equal invalid zip codes",
			a:    "22222222",
			b:    "22222222",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ca, _ := entity.NewCep(tt.a)
			cb, _ := entity.NewCep(tt.b)

			gotA := ca.Equal(cb)
			gotB := cb.Equal(ca)

			require.Equal(t, tt.want, gotA)
			require.Equal(t, tt.want, gotB)
		})
	}
}

func TestCep_MarshalJSON(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		want          []byte
		expectedError string
	}{
		{
			name:  "valid zip code",
			input: "12345678",
			want:  []byte(`"12345678"`),
		},
		{
			name:  "invalid zip code",
			input: "1234567",
			want:  []byte(`null`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := entity.NewCep(tt.input)
			got, err := c.MarshalJSON()

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestCep_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name          string
		input         []byte
		want          string
		expectedError string
	}{
		{
			name:  "valid zip code",
			input: []byte(`"12345678"`),
			want:  "12345678",
		},
		{
			name:          "invalid zip code",
			input:         []byte(`"1234567"`),
			want:          "",
			expectedError: "zip code: invalid length",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var c entity.Cep
			err := c.UnmarshalJSON(tt.input)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Zero(t, c)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, c.String())
			}
		})
	}
}
