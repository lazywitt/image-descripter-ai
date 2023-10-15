package cron

import (
	"context"
	"fmt"

	"heyalley-server/db"
)

type Pipeline interface {
	Fetch(context.Context) (PipelinePacket, error)
	Update(context.Context, PipelinePacket) error
}

type ImagePipeline struct {
	QueryEngine *db.Handler
}

func (p *ImagePipeline) Fetch(ctx context.Context) (PipelinePacket, error) {
	res, err := p.QueryEngine.GetUnprocessedImage(ctx)
	if err != nil {
		return nil, err
	}
	return &ImagePacket{
		Image: res,
	}, nil
}

func (p *ImagePipeline) Update(ctx context.Context, pp PipelinePacket) error {
	ImagePipeline, ok := pp.(*ImagePacket)
	if !ok {
		return fmt.Errorf("incorrect implementation of pipeline packet")
	}
	return p.QueryEngine.UpdateImage(ImagePipeline.Image)
}
