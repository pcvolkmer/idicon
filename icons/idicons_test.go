package icons

import (
	"bytes"
	"image/png"
	"testing"
)

func TestIgnoreCase(t *testing.T) {
	w1 := bytes.NewBuffer([]byte{})
	ig1 := NewIdIconGenerator().WithColorGenerator(ColorV1)
	png.Encode(w1, ig1.GenIcon("example", 80))

	w2 := bytes.NewBuffer([]byte{})
	ig2 := NewIdIconGenerator().WithColorGenerator(ColorV1)
	png.Encode(w2, ig2.GenIcon("Example", 80))

	if bytes.Compare(w1.Bytes(), w2.Bytes()) != 0 {
		t.Errorf("resulting images do not match")
	}
}

func TestStringMatchesHash(t *testing.T) {
	w1 := bytes.NewBuffer([]byte{})
	ig1 := NewIdIconGenerator().WithColorGenerator(ColorV2)
	// MD5 of lowercase 'example'
	png.Encode(w1, ig1.GenIcon("1a79a4d60de6718e8e5b326e338ae533", 80))

	w2 := bytes.NewBuffer([]byte{})
	ig2 := NewIdIconGenerator().WithColorGenerator(ColorV2)
	png.Encode(w2, ig2.GenIcon("Example", 80))

	if bytes.Compare(w1.Bytes(), w2.Bytes()) != 0 {
		t.Errorf("resulting images do not match")
	}
}
