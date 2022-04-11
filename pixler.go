package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

func downscale(s int, img image.Image) (image.Image, error) {
	b := img.Bounds()

	if b.Dx()%s != 0 || b.Dy()%s != 0 || s < 1 {
		return nil, fmt.Errorf("Image dimensions (%v, %v) not divisible by %v", b.Dx(), b.Dy(), s)
	}

	newImg := image.NewRGBA(
		image.Rect(b.Min.X/s, b.Min.Y/s, b.Max.X/s, b.Max.Y/s),
	)

	nB := newImg.Bounds()

	for y := nB.Min.Y; y < nB.Max.Y; y++ {
		for x := nB.Min.X; x < nB.Max.X; x++ {
			newImg.Set(x, y, img.At(s*x, s*y))
		}
	}

	return newImg, nil
}

func upscale(s int, img image.Image) image.Image {
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
	//img, _ := getImg("scaled.png")
	//img = upscale(5, img)
	//img, _ = downscale(5, img)
	//saveImg(img, "unscaled.png")
}
