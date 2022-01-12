package main

import (
	"bytes"
	_ "embed"
	"image/png"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gorilla/mux"
)

//go:embed testdata/1a79a4d60de6718e8e5b326e338ae533_v1.png
var v1 []byte

//go:embed testdata/1a79a4d60de6718e8e5b326e338ae533_v2.png
var v2 []byte

//go:embed testdata/1a79a4d60de6718e8e5b326e338ae533_gh.png
var gh []byte

func testRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/avatar/{id}", RequestHandler)
	return router
}

func TestCorrectContentType(t *testing.T) {
	req, err := http.NewRequest("GET", "/avatar/example", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	testRouter().ServeHTTP(rr, req)

	if ctype := rr.Header().Get("Content-Type"); ctype != "image/png" {
		t.Errorf("content type header does not match: got %v want image/png", ctype)
	}
}

func TestCorrectResponseForV1ColorScheme(t *testing.T) {
	req, err := http.NewRequest("GET", "/avatar/1a79a4d60de6718e8e5b326e338ae533?c=v1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	testRouter().ServeHTTP(rr, req)

	if !reflect.DeepEqual(rr.Body.Bytes(), v1) {
		t.Errorf("returned image does not match expected image")
	}
}

func TestCorrectResponseForV2ColorScheme(t *testing.T) {
	req, err := http.NewRequest("GET", "/avatar/1a79a4d60de6718e8e5b326e338ae533?c=v2", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	testRouter().ServeHTTP(rr, req)

	if !reflect.DeepEqual(rr.Body.Bytes(), v2) {
		t.Errorf("returned image does not match expected image")
	}
}

func TestCorrectResponseForGHColorSchemeAndPattern(t *testing.T) {
	req, err := http.NewRequest("GET", "/avatar/1a79a4d60de6718e8e5b326e338ae533?c=gh&d=github", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	testRouter().ServeHTTP(rr, req)

	if !reflect.DeepEqual(rr.Body.Bytes(), gh) {
		t.Errorf("returned image does not match expected image")
	}
}

func TestIgnoreCase(t *testing.T) {
	w1 := bytes.NewBuffer([]byte{})
	png.Encode(w1, genIdIcon("example", 80, colorV1))

	w2 := bytes.NewBuffer([]byte{})
	png.Encode(w2, genIdIcon("Example", 80, colorV1))

	if bytes.Compare(w1.Bytes(), w2.Bytes()) != 0 {
		t.Errorf("resulting images do not match")
	}
}

func TestStringMatchesHash(t *testing.T) {
	w1 := bytes.NewBuffer([]byte{})
	// MD5 of lowercase 'example'
	png.Encode(w1, genIdIcon("1a79a4d60de6718e8e5b326e338ae533", 80, colorV2))

	w2 := bytes.NewBuffer([]byte{})
	png.Encode(w2, genIdIcon("Example", 80, colorV2))

	if bytes.Compare(w1.Bytes(), w2.Bytes()) != 0 {
		t.Errorf("resulting images do not match")
	}
}

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
