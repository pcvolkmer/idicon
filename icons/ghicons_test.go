package icons

import "testing"

func TestHSLtoRGB(t *testing.T) {
	red := hslToRgba(0, 100, 50)
	if red.R != 255 || red.G != 0 || red.B != 0 {
		t.Errorf("Color red not as required")
	}

	green := hslToRgba(120, 100, 50)
	if green.R != 0 || green.G != 255 || green.B != 0 {
		t.Errorf("Color green not as required")
	}

	blue := hslToRgba(240, 100, 50)
	if blue.R != 0 || blue.G != 0 || blue.B != 255 {
		t.Errorf("Color blue not as required")
	}
}

func TestShouldCreateNibbles(t *testing.T) {
	hash := [16]byte{}
	hash[0] = 0x12
	nibbles := nibbles(hash)

	if nibbles[0] != 0x01 || nibbles[1] != 02 {
		t.Errorf("Nibbles not extracted as expected")
	}
}
