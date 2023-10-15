package cron

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"heyalley-server/db/models"
)

type PipelinePacket interface {
	Execute(context.Context) error
	IsUpdated(context.Context) bool
	GetPacketName(context.Context) string
}

type ImagePacket struct {
	Image *models.Image
}

func (i *ImagePacket) GetPacketName(context.Context) string {
	return i.Image.Path
}

func (i *ImagePacket) IsUpdated(context.Context) bool {
	return i.Image.IsPipelineProcessed
}

func (i *ImagePacket) Execute(ctx context.Context) error {
	currDirectory, err := os.Getwd()
	if err != nil {
		return err
	}
	imagesDirectory := filepath.Join(currDirectory, "images-blob")
	scriptDirectory := filepath.Join(currDirectory, "image/pipeline") + "/model.sh"

	fmt.Println("image loc:  ", imagesDirectory+"/"+i.Image.Path)
	cmd := exec.Command("bash", scriptDirectory, imagesDirectory+"/"+i.Image.Path)
	stderr, _ := cmd.StderrPipe()

	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
	}

	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	i.Image.IsPipelineProcessed = true
	i.Image.Description = sql.NullString{String: string(stdout), Valid: true}
	return nil
}
