package service

import (
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/maxik12233/image-process-microservice/pkg/logger"
	"go.uber.org/zap"
	"golang.org/x/image/bmp"
)

var (
	errUnsupportedFileFormat   = errors.New("Unsupported file format")
	errFailedDecodingImageFile = errors.New("Failed while decoding image bytes, probably corrupted file")
	errInternal                = errors.New("Internal error occured")
)

type IEndpoints interface {
	MakeResizeEndpoint(s IService) http.HandlerFunc
	MakeFileUploadEndpoint(s IService) http.HandlerFunc
}

type Endpoints struct {
	service IService
}

type CommonResponse struct {
	Text string `json:"text"`
	Code int    `json:"code"`
}

func NewEndpoints(s IService) IEndpoints {
	return &Endpoints{
		service: s,
	}
}

func (e *Endpoints) MakeFileUploadEndpoint(s IService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.GetLogger()

		file, _, err := r.FormFile("file")
		if err != nil {
			log.Error("Error while parsing file from form", zap.Error(err))
			return
		}

		defer file.Close()

		out, err := os.Create("uploadedfile")
		if err != nil {
			log.Error("Unable to create output file", zap.Error(err))
			return
		}

		defer out.Close()

		_, err = io.Copy(out, file)
		if err != nil {
			log.Error("Failed copying the file", zap.Error(err))
			return
		}

		log.Info("File uploaded sucesessfully", zap.Error(err))
	}
}

func (e *Endpoints) MakeResizeEndpoint(s IService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.GetLogger()

		file, header, err := r.FormFile("file")

		if err != nil {
			log.Error("Error while parsing file from form", zap.Error(err))
			return
		}

		var img image.Image
		var resimg image.Image

		filestrings := strings.Split(header.Filename, ".")
		fileformat := filestrings[len(filestrings)-1]
		switch fileformat {
		case "bmp":
			img, err = bmp.Decode(file)
			if err != nil {
				log.Error("Failed decoding image file", zap.Error(err))
				writeError(w, http.StatusBadRequest, errFailedDecodingImageFile.Error())
				return
			}
		case "jpeg", "jpg":
			img, err = jpeg.Decode(file)
			if err != nil {
				log.Error("Failed decoding image file", zap.Error(err))
				writeError(w, http.StatusBadRequest, errFailedDecodingImageFile.Error())
				return
			}
		default:
			log.Error("Unsupported file format")
			writeError(w, http.StatusBadRequest, errUnsupportedFileFormat.Error())
			return
		}

		resimg, err = s.ResizeImage(img)
		if err != nil {
			log.Error("Failed processing image file", zap.Error(err))
			writeError(w, http.StatusBadGateway, errInternal.Error())
			return
		}

		contenttype := fmt.Sprintf("image/%s", fileformat)
		disposition := fmt.Sprintf(`attachment;filename="%s"`, "resized"+header.Filename)
		w.Header().Set("Content-Type", contenttype)
		w.Header().Set("Content-Disposition", disposition)

		switch fileformat {
		case "bmp":
			if err := bmp.Encode(w, resimg); err != nil {
				log.Error("Failed encoding image file", zap.Error(err))
				writeError(w, http.StatusBadGateway, errInternal.Error())
				return
			}
		case "jpeg", "jpg":
			if err := jpeg.Encode(w, resimg, nil); err != nil {
				log.Error("Failed encoding image file", zap.Error(err))
				writeError(w, http.StatusBadGateway, errInternal.Error())
				return
			}
		}

		w.WriteHeader(http.StatusOK)
	}
}

func writeError(w http.ResponseWriter, status int, errMsg string) {
	w.WriteHeader(status)
	w.Write([]byte(errMsg))
}
