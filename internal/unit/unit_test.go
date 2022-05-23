package unit

import (
	"errors"
	"testing"
)

func TestParseSize(t *testing.T) {
	type args struct {
		inputSize         string
		expectedFloatSize float64
		expectedUnit      Unit
		err               error
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "valid int size giga",
			args: args{
				inputSize:         "2Gi",
				expectedFloatSize: 2.0,
				expectedUnit:      giga,
				err:               nil,
			},
		},
		{
			name: "valid float size giga",
			args: args{
				inputSize:         "2.2Gi",
				expectedFloatSize: 2.2,
				expectedUnit:      giga,
				err:               nil,
			},
		},
		{
			name: "valid int size mega",
			args: args{
				inputSize:         "500Mi",
				expectedFloatSize: 500,
				expectedUnit:      mega,
				err:               nil,
			},
		},
		{
			name: "invalid size, raise error",
			args: args{
				inputSize:         "2GB",
				expectedFloatSize: 0,
				expectedUnit:      unsupported,
				err:               errors.New("invalid size"),
			},
		},
		{
			name: "invalid size, raise error",
			args: args{
				inputSize:         "2.2.2GB",
				expectedFloatSize: 0,
				expectedUnit:      unsupported,
				err:               errors.New("invalid size"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, u, err := ParseSize(tt.args.inputSize)
			if err != nil && tt.args.err == nil {
				t.Errorf("error raised -> %v", err)
			}
			if tt.args.err != nil && err == nil {
				t.Errorf("failed to return error")
			}
			if tt.args.err != nil && err != nil {
				if tt.args.err.Error() != err.Error() {
					t.Errorf("wrong error returned -> got: %v, want: %v", tt.args.err.Error(), err.Error())
				}
			}
			if f != tt.args.expectedFloatSize {
				t.Errorf("invalid parsed float -> got: %v, want: %v", f, tt.args.expectedFloatSize)
			}
			if u != tt.args.expectedUnit {
				t.Errorf("invalid unit -> got: %v, want %v", u, tt.args.expectedUnit)
			}
		})
	}

}
