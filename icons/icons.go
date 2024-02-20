package icons

import (
	"crypto/md5"
	"encoding/hex"
	"image"
	"image/color"
	"image/draw"
	"regexp"
	"strings"
)

type IconGenerator interface {
	GenIcon(id string, size int) *image.NRGBA
}

func HashBytes(id string) [16]byte {
	hash := [16]byte{}
	md5RegExp := regexp.MustCompile("[a-f0-9]{32}")
	if len(id) == 32 && md5RegExp.MatchString(strings.ToLower(id)) {
		dec, _ := hex.DecodeString(id)
		for idx, b := range dec {
			print(idx)
			hash[idx] = b
		}
	} else {
		hash = md5.Sum([]byte(strings.ToLower(id)))
	}
	return hash
}

func mirrorData(data []bool, blocks int) []bool {
	for x := 0; x < blocks; x++ {
		min := x*blocks + 1
		for y := 0; y < blocks; y++ {
			a := ((blocks - x - 1) * blocks) + y
			b := min + y - 1
			if data[a] {
				data[b] = true
			}
		}
	}
	return data
}

func drawImage(data []bool, blocks int, size int, c color.Color) *image.NRGBA {
	img := image.NewNRGBA(image.Rect(0, 0, size, size))

	draw.Draw(img, img.Bounds(), &image.Uniform{color.Gray{240}}, image.Point{0, 0}, draw.Src)

	blockSize := size / (blocks + 1)
	border := (size - (blocks * blockSize)) / 2

	for x := border; x < blockSize*blocks+border; x++ {
		bx := (x - border) / blockSize
		for y := border; y < blockSize*blocks+border; y++ {
			by := (y - border) / blockSize
			idx := bx*blocks + by
			if data[idx] && (bx < blocks || by < blocks) {
				img.Set(x, y, c)
			}
		}
	}

	return img
}
