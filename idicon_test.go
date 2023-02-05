package main

import (
	_ "embed"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/gorilla/mux"
)

//go:embed testdata/1a79a4d60de6718e8e5b326e338ae533_v1.png
var v1 []byte

//go:embed testdata/1a79a4d60de6718e8e5b326e338ae533_v2.png
var v2 []byte

//go:embed testdata/1a79a4d60de6718e8e5b326e338ae533_gh.png
var gh1 []byte

//go:embed testdata/a1d0c6e83f027327d8461063f4ac58a6_gh.png
var gh2 []byte

func testRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/avatar/{id}", requestHandler)
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

	if !reflect.DeepEqual(rr.Body.Bytes(), gh1) {
		t.Errorf("returned image does not match expected image")
	}
}

func TestCorrectResponseForAltGHColorSchemeAndPattern(t *testing.T) {
	req, err := http.NewRequest("GET", "/avatar/1a79a4d60de6718e8e5b326e338ae533?c=github&d=gh", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	testRouter().ServeHTTP(rr, req)

	if !reflect.DeepEqual(rr.Body.Bytes(), gh1) {
		t.Errorf("returned image does not match expected image")
	}
}

func TestUsesConfig(t *testing.T) {
	configure("./testdata/testconfig.toml")

	if config.Defaults.ColorScheme != "gh" ||
		config.Users[0].ID != "example" ||
		config.Users[0].Alias != "42" ||
		config.Users[0].ColorScheme != "gh" ||
		config.Users[0].Pattern != "github" {
		t.Errorf("Config not applied as expected")
	}
}

func TestUsesConfigWithEnvVar(t *testing.T) {
	os.Setenv("COLORSCHEME", "v1")
	os.Setenv("PATTERN", "default")

	configure("./testdata/testconfig.toml")

	if config.Defaults.ColorScheme != "v1" ||
		config.Users[0].ID != "example" ||
		config.Users[0].Alias != "42" ||
		config.Users[0].ColorScheme != "gh" ||
		config.Users[0].Pattern != "github" {
		t.Errorf("Config not applied as expected")
	}
}

func TestCorrectResponseForUserConfig(t *testing.T) {
	configure("./testdata/testconfig.toml")

	req, err := http.NewRequest("GET", "/avatar/example", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	testRouter().ServeHTTP(rr, req)

	if !reflect.DeepEqual(rr.Body.Bytes(), gh2) {
		t.Errorf("returned image does not match expected image for mapped alias '42'")
	}
}
