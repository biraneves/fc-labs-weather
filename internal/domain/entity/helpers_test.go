package entity_test

import (
	"math"
	"testing"

	"github.com/biraneves/fc-labs-weather/internal/domain/entity"
	"github.com/stretchr/testify/require"
)

func TestAlmostEqual(t *testing.T) {
	tests := []struct {
		name    string
		a       float64
		b       float64
		epsilon float64
		want    bool
	}{
		{
			name:    "equal numbers",
			a:       25.0,
			b:       25.0,
			epsilon: 0.0,
			want:    true,
		},
		{
			name:    "difference is less than epsilon",
			a:       10.0,
			b:       10.0005,
			epsilon: 0.001,
			want:    true,
		},
		{
			name:    "difference is equal epsilon",
			a:       -3.2,
			b:       -3.1,
			epsilon: 0.10000000000000009,
			want:    true,
		},
		{
			name:    "difference greater than epsilon",
			a:       0.0,
			b:       0.002,
			epsilon: 0.001,
			want:    false,
		},
		{
			name:    "negative values inside the limit",
			a:       -15.25,
			b:       -15.249,
			epsilon: 0.002,
			want:    true,
		},
		{
			name:    "comparison with nan",
			a:       math.NaN(),
			b:       1.0,
			epsilon: 0.1,
			want:    false,
		},
		{
			name:    "comparison with infinite",
			a:       math.Inf(1),
			b:       math.Inf(1),
			epsilon: 0.1,
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := entity.AlmostEqual(tt.a, tt.b, tt.epsilon)
			require.Equal(t, tt.want, got)
		})
	}
}
