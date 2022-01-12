package main

import "image/color"

func colorV1(hash [16]byte) color.RGBA {
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

func colorV2(hash [16]byte) color.RGBA {
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

func colorGh(hash [16]byte) color.RGBA {
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
		uint8(hueToRgb(t1, t2, hue+1.0/3.0) * 255),
		uint8(hueToRgb(t1, t2, hue) * 255),
		uint8(hueToRgb(t1, t2, hue-1.0/3.0) * 255),
		0xff,
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
