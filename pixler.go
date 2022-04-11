package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

func scaleDown(s int, img image.Image) image.Image {
	return nil
}

func scaleUp(s int, img image.Image) image.Image {
	b := img.Bounds()

	newImg := image.NewRGBA(
		image.Rect(s*b.Min.X, s*b.Min.Y, s*b.Max.X, s*b.Max.Y),
	)

	for y := b.Min.Y; y < b.Max.Y; y++ {
		sY := s * y
		for i := 0; i < s; i++ {
			for x := b.Min.X; x < b.Max.X; x++ {
				sX := s * x
				for j := 0; j < s; j++ {
					newImg.Set(sX+j, sY+i, img.At(x, y))
				}
			}
		}
	}

	return newImg
}

func getImg(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	return img, err
}

func saveImg(img image.Image, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	switch ext := strings.ToLower(filepath.Ext(path)[1:]); ext {
	case "png":
		png.Encode(f, img)
	default:
		return fmt.Errorf("Unsupported filetype: %v", ext)
	}

	return nil
}

func main() {
	//img, _ := getImg("test.png")
	//img = scaleUp(5, img)
	//saveImg(img, "scaled.png")
}
