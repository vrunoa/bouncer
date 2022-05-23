package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"io"
	"testing"
)

type MockCommonApiClient struct {
	ImageListFn func(ctx context.Context, options types.ImageListOptions) ([]types.ImageSummary, error)
	ImagePullFN func(ctx context.Context, ref string, options types.ImagePullOptions) (io.ReadCloser, error)
}

func (m *MockCommonApiClient) ImageList(ctx context.Context, options types.ImageListOptions) ([]types.ImageSummary, error) {
	return m.ImageListFn(ctx, options)
}

func (m *MockCommonApiClient) ImagePull(ctx context.Context, ref string, options types.ImagePullOptions) (io.ReadCloser, error) {
	return m.ImagePullFN(ctx, ref, options)
}

func TestHandler_HasImage(t *testing.T) {
	h := &Handler{
		client: &MockCommonApiClient{
			ImageListFn: func(ctx context.Context, options types.ImageListOptions) ([]types.ImageSummary, error) {
				img := types.ImageSummary{
					ID: "whatever",
				}
				return []types.ImageSummary{img}, nil
			},
		},
	}
	ok, err := h.HasImage(context.Background(), "myimage")
	if err != nil {
		t.Errorf("error raised -> %v", err)
	}
	if !ok {
		t.Errorf("should have found an image")
	}
}

func TestHandler_HasImage_NotFound(t *testing.T) {
	h := &Handler{
		client: &MockCommonApiClient{
			ImageListFn: func(ctx context.Context, options types.ImageListOptions) ([]types.ImageSummary, error) {
				return []types.ImageSummary{}, nil
			},
		},
	}
	ok, err := h.HasImage(context.Background(), "myimage")
	if err != nil {
		t.Errorf("error raised -> %v", err)
	}
	if ok {
		t.Errorf("should have not found an image")
	}
}
