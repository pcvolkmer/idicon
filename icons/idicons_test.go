package icons

import (
	"bytes"
	"image/png"
	"testing"
)

func TestIgnoreCase(t *testing.T) {
	iconGenerator := IdIconGenerator{}

	w1 := bytes.NewBuffer([]byte{})
	png.Encode(w1, iconGenerator.GenIcon("example", 80, ColorV1))

	w2 := bytes.NewBuffer([]byte{})
	png.Encode(w2, iconGenerator.GenIcon("Example", 80, ColorV1))

	if bytes.Compare(w1.Bytes(), w2.Bytes()) != 0 {
		t.Errorf("resulting images do not match")
	}
}

func TestStringMatchesHash(t *testing.T) {
	iconGenerator := IdIconGenerator{}

	w1 := bytes.NewBuffer([]byte{})
	// MD5 of lowercase 'example'
	png.Encode(w1, iconGenerator.GenIcon("1a79a4d60de6718e8e5b326e338ae533", 80, ColorV2))

	w2 := bytes.NewBuffer([]byte{})
	png.Encode(w2, iconGenerator.GenIcon("Example", 80, ColorV2))

	if bytes.Compare(w1.Bytes(), w2.Bytes()) != 0 {
		t.Errorf("resulting images do not match")
	}
}
