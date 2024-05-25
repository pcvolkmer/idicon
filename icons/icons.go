package icons

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"regexp"
	"strings"
)

type IconGenerator interface {
	GenIcon(id string, size int) *image.NRGBA
	GenSvg(id string, size int) string
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
		minBlock := x*blocks + 1
		for y := 0; y < blocks; y++ {
			a := ((blocks - x - 1) * blocks) + y
			b := minBlock + y - 1
			if data[a] {
				data[b] = true
			}
		}
	}
	return data
}

func drawImage(data []bool, blocks int, size int, c color.Color) *image.NRGBA {
	img := image.NewNRGBA(image.Rect(0, 0, size, size))

	draw.Draw(img, img.Bounds(), &image.Uniform{C: color.Gray{Y: 240}}, image.Point{X: 0, Y: 0}, draw.Src)

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

func drawSvg(data []bool, blocks int, size int, c color.Color) string {

	blockSize := size / (blocks + 1)
	border := (size - (blocks * blockSize)) / 2
	r, g, b, _ := c.RGBA()
	colorHtml := fmt.Sprintf("#%x%x%x", r>>8, g>>8, b>>8)

	blockElems := fmt.Sprintf("<rect style=\"fill:#f0f0f0\" width=\"%d\" height=\"%d\" x=\"0\" y=\"0\" />", size, size)

	for x := 0; x < blocks; x++ {
		for y := 0; y < blocks; y++ {
			idx := x*blocks + y
			if data[idx] {
				blockElems += fmt.Sprintf(
					`<rect style="fill:%s" width="%d" height="%d" x="%d" y="%d" />`,
					colorHtml,
					blockSize,
					blockSize,
					border+(x*blockSize),
					border+(y*blockSize))
			}
		}
	}

	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<svg width="%d" height="%d" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:svg="http://www.w3.org/2000/svg"><g>%s</g></svg>`,
		size, size,
		blockElems)
}
