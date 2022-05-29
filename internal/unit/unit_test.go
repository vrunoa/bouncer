package unit

import (
	"testing"
)

func TestToBytes(t *testing.T) {
	type args struct {
		inputSize    float64
		inputUnit    Unit
		expectedSize int64
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "check mega to bytes",
			args: args{
				inputSize:    1.0,
				inputUnit:    Mega,
				expectedSize: 1036288,
			},
		},
		{
			name: "check giga to bytes",
			args: args{
				inputSize:    1.0,
				inputUnit:    Giga,
				expectedSize: 1061158912,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := ToBytes(tt.args.inputSize, tt.args.inputUnit)
			if res != tt.args.expectedSize {
				t.Errorf("wrong conversion. Want: %v Got: %v", tt.args.expectedSize, res)
			}
		})
	}
}
