package service

import (
	"image"

	"github.com/maxik12233/image-process-microservice/pkg/pixels"
)

type IService interface {
	ResizeImage(img image.Image) (image.Image, error)
}

type Service struct {
}

func NewService() IService {
	return &Service{}
}

func (s *Service) ResizeImage(img image.Image) (image.Image, error) {

	newimg := pixels.NewPixelImage(img)

	err := newimg.ToFile("image", false)
	if err != nil {
		return nil, err
	}

	return newimg.GetImage(), nil
}
