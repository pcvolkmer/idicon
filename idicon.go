package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/gorilla/mux"
)

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

// Based on https://github.com/dgraham/identicon
func genGhIcon(id string, size int, f func([16]byte) color.RGBA) *image.NRGBA {
	if size > 512 {
		size = 512
	}
	blocks := 5

	hash := hashBytes(id)
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

	return drawImage(mirrorData(data, blocks), blocks, size, f(hash))
}

func genIdIcon(id string, size int, f func([16]byte) color.RGBA) *image.NRGBA {
	id = strings.ToLower(id)
	blocks := 5
	if size > 512 {
		size = 512
	}

	hash := hashBytes(id)
	data := make([]bool, blocks*blocks)
	for i := 0; i < len(hash)-1; i++ {
		data[i] = hash[i]%2 != hash[i+1]%2
	}
	return drawImage(mirrorData(data, blocks), blocks, size, f(hash))
}

func hashBytes(id string) [16]byte {
	hash := [16]byte{}
	md5RegExp := regexp.MustCompile("[a-f0-9]{32}")
	if !md5RegExp.Match([]byte(id)) {
		hash = md5.Sum([]byte(id))
	} else {
		dec, _ := hex.DecodeString(id)
		for idx, b := range dec {
			hash[idx] = b
		}
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

func requestHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	size, err := strconv.Atoi(r.URL.Query().Get("s"))
	if err != nil {
		size = 80
	}

	colorScheme := r.URL.Query().Get("c")
	if colorScheme == "" {
		colorScheme = config.Defaults.ColorScheme
	}

	pattern := r.URL.Query().Get("d")
	if pattern == "" {
		pattern = config.Defaults.Pattern
	}

	for _, userConfig := range config.Users {
		if hashBytes(id) == hashBytes(userConfig.ID) {
			id = userConfig.Alias
			if len(userConfig.ColorScheme) > 0 {
				colorScheme = userConfig.ColorScheme
			}
			if len(userConfig.Pattern) > 0 {
				pattern = userConfig.Pattern
			}
		}
	}

	w.Header().Add("Content-Type", "image/png")
	cFunc := colorV2
	if colorScheme == "v1" {
		cFunc = colorV1
	} else if colorScheme == "gh" {
		cFunc = colorGh
	}

	if pattern == "github" {
		err = png.Encode(w, genGhIcon(id, size, cFunc))
	} else {
		err = png.Encode(w, genIdIcon(id, size, cFunc))
	}

}

var (
	config Config
)

func configure(configFile string) {
	if file, err := os.OpenFile(configFile, os.O_RDONLY, 0); err == nil {
		c := &Config{}
		_, err := toml.NewDecoder(file).Decode(c)
		if err != nil {
			log.Printf("Invalid config file '%s' - ignore it.\n", configFile)
		}

		defer file.Close()
		config = *c
	}

	if os.Getenv("COLORSCHEME") != "" {
		config.Defaults.ColorScheme = os.Getenv("COLORSCHEME")
	}

	if os.Getenv("PATTERN") != "" {
		config.Defaults.Pattern = os.Getenv("PATTERN")
	}
}

func main() {
	configFile := flag.String("c", "/etc/idicon/config.toml", "-c <path to config file>")
	flag.Parse()

	configure(*configFile)

	router := mux.NewRouter()
	router.HandleFunc("/avatar/{id}", requestHandler)
	log.Println("Starting ...")
	log.Fatal(http.ListenAndServe(":8000", router))
}
