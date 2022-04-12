package main

import (
	"encoding/hex"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func Downscale(s int, img image.Image) (image.Image, error) {
	b := img.Bounds()

	if s < 1 {
		return nil, fmt.Errorf("Invalid value of s: %v", s)
	}

	/*if b.Dx()%s != 0 || b.Dy()%s != 0 {
		return nil, fmt.Errorf("Image dimensions (%v, %v) not divisible by %v", b.Dx(), b.Dy(), s)
	}*/

	newImg := image.NewRGBA(
		image.Rect(b.Min.X/s, b.Min.Y/s, b.Max.X/s, b.Max.Y/s),
	)

	nB := newImg.Bounds()
	o := s / 2

	for y := nB.Min.Y; y < nB.Max.Y; y++ {
		for x := nB.Min.X; x < nB.Max.X; x++ {
			newImg.Set(x, y, img.At(s*x+o, s*y+o))
		}
	}

	return newImg, nil
}

func Upscale(s int, img image.Image) image.Image {
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

func colorEq(a, b color.Color) bool {
	aR, aG, aB, aA := a.RGBA()
	bR, bG, bB, bA := b.RGBA()

	return (aR == bR &&
		aG == bG &&
		aB == bB &&
		aA == bA)
}

func rgbEq(a, b color.Color) bool {
	aR, aG, aB, _ := a.RGBA()
	bR, bG, bB, _ := b.RGBA()

	return (aR == bR &&
		aG == bG &&
		aB == bB)
}

func pxsEq(img image.Image, r image.Rectangle, o image.Point) bool {
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			cc := img.At(x, y)
			nc := img.At(x+o.X, y+o.Y)
			if !colorEq(cc, nc) {
				return false
			}
		}
	}

	return true
}

func Unrepeat(img image.Image) image.Image {
	b := img.Bounds()
	tempImg := image.NewRGBA(b)
	j := 0
	for y := b.Min.Y; y < b.Max.Y; y++ {
		ln := image.Rect(b.Min.X, y, b.Max.X, y+1)
		if pxsEq(img, ln, image.Pt(0, 1)) {
			continue
		} else {
			cl := image.Rect(b.Min.X, j, b.Max.X, j+1)
			draw.Draw(tempImg, cl, img, image.Pt(0, y), draw.Src)
			j++
		}
	}

	tB := image.Rect(0, 0, b.Max.X, j)
	i := 0
	for x := b.Min.X; x < b.Max.X; x++ {
		ln := image.Rect(x, b.Min.Y, x+1, tB.Max.Y)
		if pxsEq(tempImg, ln, image.Pt(1, 0)) {
			continue
		} else {
			cl := image.Rect(i, b.Min.Y, i+1, tB.Max.Y)
			draw.Draw(tempImg, cl, tempImg, image.Pt(x, 0), draw.Src)
			i++
		}
	}

	nB := image.Rect(0, 0, i, j)
	newImg := image.NewRGBA(nB)
	draw.Draw(newImg, nB, tempImg, nB.Min, draw.Src)

	return newImg
}

func loadPalette(path string) (color.Palette, error) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	lns := strings.Split(strings.Replace(string(f), "\r\n", "\n", -1), "\n")
	var pl color.Palette
	for _, hexStr := range lns {
		if hexStr == "" {
			continue
		}
		bc, err := hex.DecodeString(hexStr[1:])
		if err != nil {
			return nil, err
		}
		r := bc[0]
		g := bc[1]
		b := bc[2]
		var a uint8
		if len(bc) > 3 {
			a = bc[3]
		} else {
			a = 255
		}
		pl = append(pl, color.RGBA{r, g, b, a})
	}
	return pl, nil
}

func PaletteDiff(img image.Image, p color.Palette) (image.Image, int) {
	b := img.Bounds()

	diff := image.NewRGBA(b)
	n := 0

	//fmt.Println(p)

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			valid := false
			for _, c := range p {
				if rgbEq(c, img.At(x, y)) {
					valid = true
					break
				}
			}
			if !valid {
				//fmt.Printf("%v,%v: %v | %v\n", x, y, c, img.At(x, y))
				diff.Set(x, y, img.At(x, y))
				n++
			}
		}
	}

	return diff, n
}

// TODO: all these funcs are redundant, clean up

func parseScaleArg(i int) int {
	scale, err := strconv.Atoi(os.Args[i])
	if err != nil {
		log.Fatal(err)
	}
	return scale
}

func getImgArg(i int) image.Image {
	in := os.Args[i]
	img, err := getImg(in)
	if err != nil {
		log.Fatal(err)
	}

	return img
}

func getOutArg(i int) string {
	if len(os.Args) > i {
		return os.Args[i]
	} else {
		return "pixler-output.png"
	}
}

func getStrArg(i int, dflt string) string {
	if len(os.Args) > i {
		return os.Args[i]
	} else {
		return dflt
	}
}

func main() {
	if len(os.Args) < 2 {
		os.Exit(2)
	}
	cmd := os.Args[1]
	var err error
	switch cmd {
	case "upscale":
		scale := parseScaleArg(2)
		img := getImgArg(3)
		out := getOutArg(4)
		img = Upscale(scale, img)
		err = saveImg(img, out)
	case "downscale":
		scale := parseScaleArg(2)
		img := getImgArg(3)
		out := getOutArg(4)
		img, err = Downscale(scale, img)
		if err != nil {
			log.Fatal(err)
		}
		err = saveImg(img, out)
	case "unrepeat":
		img := getImgArg(2)
		out := getOutArg(3)
		img = Unrepeat(img)
		err = saveImg(img, out)
	case "palette":
		if len(os.Args) < 3 {
			panic("not enough arguments to command: palette")
		}
		img := getImgArg(2)
		out := getStrArg(3, "diff.png")
		pl, err := loadPalette(getStrArg(4, "palette.txt"))
		if err != nil {
			log.Fatal(err)
		}
		diff, n := PaletteDiff(img, pl)
		saveImg(diff, out)
		fmt.Printf("%v pixels differ from the palette\n", n)
	default:
		log.Fatalf("Unknown command '%v'", cmd)
	}
	if err != nil {
		log.Fatal(err)
	}
}
