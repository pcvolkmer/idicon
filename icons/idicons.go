package icons

import (
	"image"
	"image/color"
	"strings"
)

type IdIconGenerator struct {
}

func (generator *IdIconGenerator) GenIcon(id string, size int, f func([16]byte) color.RGBA) *image.NRGBA {
	id = strings.ToLower(id)
	blocks := 5
	if size > 512 {
		size = 512
	}

	hash := HashBytes(id)
	data := make([]bool, blocks*blocks)
	for i := 0; i < len(hash)-1; i++ {
		data[i] = hash[i]%2 != hash[i+1]%2
	}
	return drawImage(mirrorData(data, blocks), blocks, size, f(hash))
}

func ColorV1(hash [16]byte) color.RGBA {
	r := 32 + (hash[0]%16)/2<<4
	g := 32 + (hash[2]%16)/2<<4
	b := 32 + (hash[len(hash)-1]%16)/2<<4

	if r > g && r > b {
		r += 48
	} else if g > r && g > b {
		g += 48
	} else if b > r && b > g {
		b += 48
	}
	return color.RGBA{r, g, b, 255}
}

func ColorV2(hash [16]byte) color.RGBA {
	var palette = []color.RGBA{
		{0x3c, 0x38, 0x36, 0xff},
		{0xcc, 0x24, 0x1d, 0xff},
		{0x98, 0x97, 0x1a, 0xff},
		{0xd7, 0x99, 0x21, 0xff},
		{0x45, 0x85, 0x88, 0xff},
		{0xb1, 0x62, 0x86, 0xff},
		{0x68, 0x9d, 0x6a, 0xff},
		{0xa8, 0x99, 0x84, 0xff},
	}
	return palette[hash[15]%8]
}
