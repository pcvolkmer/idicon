package main

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/gorilla/mux"
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
)

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

func RequestHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	size, err := strconv.Atoi(r.URL.Query().Get("s"))
	if err != nil {
		size = 80
	}

	colorScheme := r.URL.Query().Get("c")
	if colorScheme == "" {
		colorScheme = os.Getenv("COLORSCHEME")
	}

	w.Header().Add("Content-Type", "image/png")
	if colorScheme == "v1" {
		err = png.Encode(w, genIdIcon(id, size, colorV1))
	} else {
		err = png.Encode(w, genIdIcon(id, size, colorV2))
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/avatar/{id}", RequestHandler)
	log.Fatal(http.ListenAndServe(":8000", router))
}
