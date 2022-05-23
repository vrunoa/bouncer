package config

import (
	"errors"
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
					Image: DockerImage{
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
					Image: DockerImage{
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
