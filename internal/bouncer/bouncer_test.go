package bouncer

import (
	"github.com/vrunoa/bouncer/internal/config"
	"github.com/vrunoa/bouncer/internal/docker"
	"github.com/vrunoa/bouncer/internal/unit"
	"testing"
)

func TestPolicyResult_String(t *testing.T) {
	pol := PolicyResult{
		Message: "we ducked up",
		Desc:    "some desc",
	}
	want := "Policy: some desc. Result: we ducked up"
	if pol.String() != want {
		t.Errorf("wrong message. Want: %v Got: %v", pol.String(), want)
	}
}

func TestCheckPolicies(t *testing.T) {
	b := &bouncer{}
	type args struct {
		policies []config.DenyPolicy
		img      docker.Image
		result   PolicyResult
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test policy success",
			args: args{
				policies: []config.DenyPolicy{
					{Size: "2.2Gi", Desc: "big image", Unit: unit.Giga, FloatSize: 2},
				},
				img: docker.Image{
					Size: 100000,
					Name: "some name",
				},
				result: PolicyResult{
					Status:  0,
					Message: "policy ok",
				},
			},
		},
		{
			name: "test policy fail",
			args: args{
				policies: []config.DenyPolicy{
					{Size: "1Mi", Desc: "really tiny image", Unit: unit.Mega, FloatSize: 1},
				},
				img: docker.Image{
					Size: 2000000,
					Name: "some name",
				},
				result: PolicyResult{
					Status:  1,
					Message: "policy failed -> image size: 2000000 - deny size: 1036288",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := b.checkPolicies(&tt.args.img, tt.args.policies)
			if err != nil {
				t.Errorf("error raised -> %v", err)
			}
			if len(results) != 1 {
				t.Errorf("wrong results. %v", results)
			}
			res := results[0]
			if res.Status != tt.args.result.Status {
				t.Errorf("wrong result Status. Want: %v Got: %v", tt.args.result.Status, res.Status)
			}
			if res.Message != tt.args.result.Message {
				t.Errorf("wrong result Message. Want: %v. Got: %v", tt.args.result.Message, res.Message)
			}
		})
	}
}
