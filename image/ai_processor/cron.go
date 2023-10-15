package cron

import (
	"context"
	"fmt"
	"time"
)

// Cron executes upto infinity in a cycle
// calls pipeline to return an Executable packet
type Cron interface {
	Trigger(context.Context) error
}

// type Pipeline interface {
// 	Fetch(context.Context) (PipelinePacket, error)
// 	Update(context.Context, PipelinePacket) error
// }

// type Processor interface {
// 	Execute(context.Context) (PipelinePacket, error)
// }

// type PipelinePacket interface {
// 	Execute(context.Context) error
// 	IsUpdated(context.Context) bool
// 	GetPacketName(context.Context) string
// }

type ImageProcessingCron struct {
	SleepTime time.Duration
	Pipeline  Pipeline
}

func (c *ImageProcessingCron) Trigger(ctx context.Context) error {
	for {
		time.Sleep(c.SleepTime)
		fmt.Print("cron triggered")
		pipelinePacket, err := c.Pipeline.Fetch(ctx)
		if err != nil {
			fmt.Println("cron error fetching: ", err)
			continue
		} else if pipelinePacket == nil {
			continue
		}
		fmt.Println("packet fetched: ", pipelinePacket.GetPacketName(ctx))

		err = pipelinePacket.Execute(ctx)
		if err != nil {
			fmt.Println("cron error execution: ", err)
		}
		if pipelinePacket.IsUpdated(ctx) {
			err = c.Pipeline.Update(ctx, pipelinePacket)
			if err != nil {
				fmt.Println("cron error updating: ", err)
				continue
			}
			fmt.Println("packet updated: ", pipelinePacket.GetPacketName(ctx))
		}
	}
}

// type ImagePipeline struct {
// 	QueryEngine *db.Handler
// }

// func (p *ImagePipeline) Fetch(ctx context.Context) (PipelinePacket, error) {
// 	res, err := p.QueryEngine.GetUnprocessedImage(ctx)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &ImagePacket{
// 		Image: res,
// 	}, nil
// }

// func (p *ImagePipeline) Update(ctx context.Context, pp PipelinePacket) error {
// 	ImagePipeline, ok := pp.(*ImagePacket)
// 	if !ok{
// 		return fmt.Errorf("incorrect implementation of pipeline packet")
// 	}
// 	return p.QueryEngine.UpdateImage(ImagePipeline.Image)
// }

// // type ImagePacket struct {
// // 	Image *models.Image
// // }

// func (i *ImagePacket) GetPacketName(context.Context) string {
// 	return i.Image.Path
// }

// func (i *ImagePacket) IsUpdated(context.Context) bool {
// 	return i.Image.IsPipelineProcessed
// }

// func (i *ImagePacket) Execute(ctx context.Context) error {
// 	currDirectory, err := os.Getwd()
// 	if err != nil {
// 		return err
// 	}
// 	imagesDirectory := filepath.Join(currDirectory, "images-blob")
// 	scriptDirectory := filepath.Join(currDirectory, "image/pipeline") + "/model.sh"

// 	fmt.Println("image loc:  ", imagesDirectory+ "/"+ i.Image.Path)
// 	cmd := exec.Command("bash", scriptDirectory, imagesDirectory+ "/" + i.Image.Path)
// 	stderr, _ := cmd.StderrPipe()

// 	stdout, err := cmd.Output()
// 	if err != nil {
// 		fmt.Println(err.Error())
// 	}

// 	scanner := bufio.NewScanner(stderr)
// 	for scanner.Scan() {
// 		fmt.Println(scanner.Text())
// 	}
// 	i.Image.IsPipelineProcessed = true
// 	i.Image.Description = sql.NullString{String: string(stdout), Valid: true}
// 	return nil
// }
