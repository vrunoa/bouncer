package config

import (
	"errors"
	"github.com/vrunoa/bouncer/internal/unit"
	"testing"
)

func TestConfiguration_Validate(t *testing.T) {
	type args struct {
		c   Configuration
		err error
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test raise on missing image name",
			args: args{
				c:   Configuration{},
				err: errors.New("missing image name"),
			},
		},
		{
			name: "test raise on missing policies",
			args: args{
				c: Configuration{
					Image: Image{
						Name: "some image",
					},
				},
				err: errors.New("missing deny policies"),
			},
		},
		{
			name: "test raise on missing wrong size value",
			args: args{
				c: Configuration{
					Image: Image{
						Name: "some image",
						Policy: GuardPolicy{
							Deny: []DenyPolicy{
								{
									Size: "2.2.2Gi",
								},
							},
						},
					},
				},
				err: errors.New("invalid size"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.c.Validate()
			if err.Error() != tt.args.err.Error() {
				t.Errorf("failed to raised proper error. Want: %v Got: %v", tt.args.err.Error(), err.Error())
			}
		})
	}
}

func TestParseSize(t *testing.T) {
	type args struct {
		inputSize         string
		expectedFloatSize float64
		expectedUnit      unit.Unit
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
				expectedUnit:      unit.Giga,
				err:               nil,
			},
		},
		{
			name: "valid float size giga",
			args: args{
				inputSize:         "2.2Gi",
				expectedFloatSize: 2.2,
				expectedUnit:      unit.Giga,
				err:               nil,
			},
		},
		{
			name: "valid int size mega",
			args: args{
				inputSize:         "500Mi",
				expectedFloatSize: 500,
				expectedUnit:      unit.Mega,
				err:               nil,
			},
		},
		{
			name: "invalid size, raise error unsupported unit",
			args: args{
				inputSize:         "2GB",
				expectedFloatSize: 0,
				expectedUnit:      unit.Unsupported,
				err:               errors.New("invalid size"),
			},
		},
		{
			name: "invalid size, raise error on wrong float",
			args: args{
				inputSize:         "2.2.2GB",
				expectedFloatSize: 0,
				expectedUnit:      unit.Unsupported,
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
