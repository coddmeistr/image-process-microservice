package pixels

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"os"
	"strings"
)

type IPixelImage interface {
	String() string
	GetImage() image.Image
	ToFile(FileNameWithPath string, alpha bool) error
}

type Pixel struct {
	r uint8
	g uint8
	b uint8
	a uint8
}

type PixelImage struct {
	img    image.RGBA
	Pixels [][]Pixel
}

func NewPixelImage(img image.Image) IPixelImage {
	imgrgba := image.NewRGBA(img.Bounds())
	draw.Draw(imgrgba, img.Bounds(), img, image.Point{}, draw.Over)
	return &PixelImage{
		img:    *imgrgba,
		Pixels: getPixels(*imgrgba),
	}
}

func (img *PixelImage) String() string {
	result := ""
	for i := 0; i < len(img.Pixels); i++ {
		for j := 0; j < len(img.Pixels[i]); j++ {
			result += fmt.Sprintf("(%v, %v, %v, %v)", img.Pixels[i][j].r, img.Pixels[i][j].g, img.Pixels[i][j].b, img.Pixels[i][j].a) + " "
		}
	}

	return result
}

func (img *PixelImage) ToFile(FileNameWithPath string, alpha bool) error {

	var maxspace int
	var template string
	if alpha {
		maxspace = len("(255, 255, 255, 255)")
		template = "(%v, %v, %v, %v)"
	} else {
		maxspace = len("(255, 255, 255)")
		template = "(%v, %v, %v)"
	}

	out, err := os.Create(FileNameWithPath)
	if err != nil {
		return err
	}

	for i := 0; i < len(img.Pixels); i++ {
		for j := 0; j < len(img.Pixels[i]); j++ {
			pixel := img.Pixels[i][j]

			var toadd string
			if alpha {
				toadd = fmt.Sprintf(template, pixel.r, pixel.g, pixel.b, pixel.a)
			} else {
				toadd = fmt.Sprintf(template, pixel.r, pixel.g, pixel.b)
			}

			if _, err := out.WriteString(stringWithFixedSpace(toadd, maxspace)); err != nil {
				return err
			}
		}
		if _, err := out.WriteString("\n"); err != nil {
			return err
		}
	}

	return nil
}

func (img *PixelImage) GetImage() image.Image {
	img.applyPixelsToImg()
	return &img.img
}

func (img *PixelImage) applyPixelsToImg() {
	for i := 0; i < len(img.Pixels); i++ {
		for j := 0; j < len(img.Pixels[i]); j++ {
			pixel := img.Pixels[i][j]
			img.img.Set(j, i, color.RGBA{pixel.r, pixel.g, pixel.b, pixel.a})
		}
	}
}

func getPixels(img image.RGBA) [][]Pixel {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	var pixels [][]Pixel
	for y := 0; y < height; y++ {
		var row []Pixel
		for x := 0; x < width; x++ {
			row = append(row, rgbaToPixel(img.At(x, y).RGBA()))
		}
		pixels = append(pixels, row)
	}

	return pixels
}

func rgbaToPixel(r uint32, g uint32, b uint32, a uint32) Pixel {
	return Pixel{uint8(r / 257), uint8(g / 257), uint8(b / 257), uint8(a / 257)}
}

func stringWithFixedSpace(str string, space int) string {
	trimmed := strings.TrimSpace(str)
	count := space - len(trimmed)
	for i := 0; i < count; i++ {
		trimmed += " "
	}
	return trimmed
}
