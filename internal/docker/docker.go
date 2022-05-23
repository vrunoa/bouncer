package docker

import (
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"io"
	"io/ioutil"
)

// CommonAPIClient interface on docker client: https://github.com/moby/moby/tree/master/client#go-client-for-the-docker-engine-api
type CommonAPIClient interface {
	ImageList(ctx context.Context, options types.ImageListOptions) ([]types.ImageSummary, error)
	ImagePull(ctx context.Context, ref string, options types.ImagePullOptions) (io.ReadCloser, error)
}

type Handler interface {
	HasImage(ctx context.Context, image string) (bool, error)
	ListImages(ctx context.Context, images []string) ([]Image, error)
	GetImageInformation(ctx context.Context, image string) (*Image, error)
	PullImage(ctx context.Context, image string) error
}

// handler handles docker API calls
type handler struct {
	client CommonAPIClient
}

// Create creates new handler
func Create() (*handler, error) {
	cl, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	h := &handler{
		client: cl,
	}
	return h, nil
}

type Image struct {
	Name string
	Size int64
}

// listImages gets a list of images from docker client using image name as reference
func (h *handler) listImages(ctx context.Context, images []string) ([]types.ImageSummary, error) {
	listFilters := filters.NewArgs()
	for _, i := range images {
		listFilters.Add("reference", i)
	}
	opts := types.ImageListOptions{
		All:     true,
		Filters: listFilters,
	}
	list, err := h.client.ImageList(ctx, opts)
	if err != nil {
		return nil, err
	}
	return list, err
}

// HasImage checks if image exists locally to docker client
func (h *handler) HasImage(ctx context.Context, image string) (bool, error) {
	list, err := h.listImages(ctx, []string{image})
	if err != nil {
		return false, err
	}
	return len(list) > 0, nil
}

// ListImages list docker images filtered by list of image names
func (h *handler) ListImages(ctx context.Context, images []string) ([]Image, error) {
	list, err := h.listImages(ctx, images)
	if err != nil {
		return nil, err
	}
	var imgList []Image
	for _, img := range list {
		imgList = append(imgList, Image{
			Name: img.ID,
			Size: img.Size,
		})
	}
	return imgList, nil
}

// GetImageInformation get docker image information
func (h *handler) GetImageInformation(ctx context.Context, image string) (*Image, error) {
	list, err := h.ListImages(ctx, []string{image})
	if err != nil {
		return &Image{}, err
	}
	if len(list) == 0 {
		return &Image{}, errors.New("image not found")
	}
	return &list[0], nil
}

// PullImage pulls docker image from registry
func (h *handler) PullImage(ctx context.Context, image string) error {
	opts := types.ImagePullOptions{}
	res, err := h.client.ImagePull(ctx, image, opts)
	if err != nil {
		return err
	}
	defer res.Close()
	stdout, err := ioutil.ReadAll(res)
	if err != nil {
		return err
	}
	fmt.Println(stdout)
	return nil
}
