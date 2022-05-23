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
	h := &handler{
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
	h := &handler{
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

func TestHandler_listImagesOptions(t *testing.T) {
	var opts types.ImageListOptions
	h := &handler{
		client: &MockCommonApiClient{
			ImageListFn: func(ctx context.Context, options types.ImageListOptions) ([]types.ImageSummary, error) {
				opts = options
				return []types.ImageSummary{}, nil
			},
		},
	}
	imageName := "animage"
	_, err := h.listImages(context.Background(), []string{imageName})
	if err != nil {
		t.Errorf("error raised -> %v", err)
	}
	if opts.All == false {
		t.Errorf("wrong ImageListOptions setup for All. Want: %v. Got: %v", true, opts.All)
	}
	var refFound = false
	for _, ref := range opts.Filters.Get("reference") {
		if ref == imageName {
			refFound = true
		}
	}
	if refFound == false {
		t.Errorf("failed to set filter")
	}
}

func TestHandler_GetImageInformation(t *testing.T) {
	imageID := "whatever"
	size := int64(200009)
	h := &handler{
		client: &MockCommonApiClient{
			ImageListFn: func(ctx context.Context, options types.ImageListOptions) ([]types.ImageSummary, error) {
				img := types.ImageSummary{
					ID:   imageID,
					Size: size,
				}
				return []types.ImageSummary{img}, nil
			},
		},
	}
	img, err := h.GetImageInformation(context.Background(), "animage")
	if err != nil {
		t.Errorf("error raised -> %v", err)
	}
	if img.Name != imageID {
		t.Errorf("wrong image.Name -> Want: %v Got: %v", imageID, img.Name)
	}
	if img.Size != size {
		t.Errorf("wrong image.Size -> Want: %v Got: %v", size, img.Size)
	}
}
