package icons

import (
	"image"
	"image/color"
)

type GhIconGenerator struct {
	colorGenerator func([16]byte) color.RGBA
}

func NewGhIconGenerator() *GhIconGenerator {
	return &GhIconGenerator{}
}

func (generator *GhIconGenerator) WithColorGenerator(colorGenerator func([16]byte) color.RGBA) *GhIconGenerator {
	generator.colorGenerator = colorGenerator
	return generator
}

// Based on https://github.com/dgraham/identicon
func (generator *GhIconGenerator) GenIcon(id string, size int) *image.NRGBA {
	if size > 512 {
		size = 512
	}
	blocks := 5

	hash := HashBytes(id)
	nibbles := nibbles(hash)
	data := make([]bool, blocks*blocks)

	for x := 0; x < blocks; x++ {
		for y := 0; y < blocks; y++ {
			ni := x + blocks*(blocks-y-1)
			if x+blocks*y > 2*blocks {
				di := (x + blocks*y) - 2*blocks
				data[di] = nibbles[ni%32]%2 == 0
			}
		}
	}

	return drawImage(mirrorData(data, blocks), blocks, size, generator.colorGenerator(hash))
}

// https://processing.org/reference/map_.html
func remap(value uint32, vmin uint32, vmax uint32, dmin uint32, dmax uint32) float32 {
	return float32((value-vmin)*(dmax-dmin)) / float32((vmax-vmin)+dmin)
}

func nibbles(hash [16]byte) []byte {
	nibbles := make([]byte, 32)

	for i := 0; i <= 15; i++ {
		nibbles[i*2+1] = hash[i] & 0x0f
		nibbles[i*2] = hash[i] & 0xf0 >> 4
	}

	return nibbles
}

func ColorGh(hash [16]byte) color.RGBA {
	h1 := (uint16(hash[12]) & 0x0f) << 8
	h2 := uint16(hash[13])

	h := uint32(h1 | h2)
	s := uint32(hash[14])
	l := uint32(hash[15])

	return hslToRgba(
		remap(h, 0, 4096, 0, 360),
		65.0-remap(s, 0, 255, 0, 20),
		75.0-remap(l, 0, 255, 0, 20),
	)
}

// http://www.w3.org/TR/css3-color/#hsl-color
func hslToRgba(hue float32, sat float32, lum float32) color.RGBA {
	hue = hue / 360.0
	sat = sat / 100.0
	lum = lum / 100.0

	t2 := lum + sat - (lum * sat)
	if lum <= 0.5 {
		t2 = lum * (sat + 1.0)
	}

	t1 := lum*2.0 - t2

	return color.RGBA{
		R: uint8(hueToRgb(t1, t2, hue+1.0/3.0) * 255),
		G: uint8(hueToRgb(t1, t2, hue) * 255),
		B: uint8(hueToRgb(t1, t2, hue-1.0/3.0) * 255),
		A: 0xff,
	}
}

func hueToRgb(t1 float32, t2 float32, hue float32) float32 {
	if hue < 0.0 {
		hue = hue + 1.0
	} else if hue >= 1.0 {
		hue = hue - 1.0
	}

	if hue < 1.0/6.0 {
		return t1 + (t2-t1)*6.0*hue
	}

	if hue < 1.0/2.0 {
		return t2
	}

	if hue < 2.0/3.0 {
		return t1 + (t2-t1)*(2.0/3.0-hue)*6.0
	}

	return t1
}
