package bouncer

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/vrunoa/bouncer/internal/config"
	"github.com/vrunoa/bouncer/internal/docker"
	"github.com/vrunoa/bouncer/internal/unit"
	"time"
)

type bouncer struct {
	Configuration config.Configuration
	DockerHandler docker.Handler
}

func NewBouncer(configFile string) (*bouncer, error) {
	bouncer := &bouncer{}
	log.Debug().Str("config-file", configFile).Msg("reading configuration file")
	cfg, err := config.ReadYaml(configFile)
	if err != nil {
		return nil, err
	}
	err = cfg.Validate()
	if err != nil {
		return nil, err
	}
	bouncer.Configuration = *cfg
	handler, err := docker.Create()
	if err != nil {
		return nil, err
	}
	bouncer.DockerHandler = handler
	return bouncer, nil
}

type PolicyResult struct {
	Status  int
	Message string
	Desc    string
}

func (p *PolicyResult) String() string {
	return fmt.Sprintf("Policy: %s. Result: %s", p.Desc, p.Message)
}

func (b *bouncer) Check() ([]PolicyResult, error) {
	ctx := context.Background()
	imageName := b.Configuration.Image.Name
	log.Debug().Str("image", imageName).Msg("checking if image is listed locally")
	has, err := b.DockerHandler.HasImage(ctx, imageName)
	if err != nil {
		return nil, err
	}
	if !has {
		log.Debug().Str("image", imageName).Msg("pulling image")
		log.Debug().Str("start-time", time.Now().String()).Msg("capturing metrics")
		err := b.DockerHandler.PullImage(ctx, imageName)
		if err != nil {
			return nil, err
		}
		log.Debug().Str("end-time", time.Now().String()).Msg("capturing metrics")
	}
	log.Debug().Str("image", imageName).Msg("getting image information")
	img, err := b.DockerHandler.GetImageInformation(ctx, imageName)
	if err != nil {
		return nil, err
	}
	return b.checkPolicies(img, b.Configuration.Image.Policy.Deny)
}

func (b *bouncer) checkPolicies(image *docker.Image, policies []config.DenyPolicy) ([]PolicyResult, error) {
	var results []PolicyResult
	for _, deny := range policies {
		log.Debug().Str("image", image.Name).Str("desc", deny.Desc).Str("size", deny.Size).Msg("checking policy")
		res := PolicyResult{
			Status:  0,
			Desc:    deny.Desc,
			Message: "policy ok",
		}
		sizeFloat, sizeUnit, err := unit.ParseSize(deny.Size)
		if err != nil {
			res.Status = 1
			res.Message = err.Error()
		} else {
			sizeBytes := unit.ToBytes(sizeFloat, sizeUnit)
			if image.Size > sizeBytes {
				res.Status = 1
				res.Message = fmt.Sprintf("policy failed -> image size: %d - deny size: %d", image.Size, sizeBytes)
			}
		}
		results = append(results, res)
	}
	return results, nil
}
