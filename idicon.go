package main

import (
	"embed"
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/gorilla/mux"
	"idicon/icons"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"net/http"
	"os"
	"strconv"
)

//go:embed static
var static embed.FS

func pageRequestHandler(w http.ResponseWriter, _ *http.Request) {
	p, _ := static.ReadFile("static/index.html")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(p)
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	size, err := strconv.Atoi(r.URL.Query().Get("s"))
	if err != nil {
		size, err = strconv.Atoi(r.URL.Query().Get("size"))
		if err != nil {
			size = 80
		}
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
		if icons.HashBytes(id) == icons.HashBytes(userConfig.ID) {
			if userConfig.Redirect != "" {
				w.Header().Add("Location", userConfig.Redirect)
				w.WriteHeader(http.StatusFound)
				return
			}

			id = userConfig.Alias
			if len(userConfig.ColorScheme) > 0 {
				colorScheme = userConfig.ColorScheme
			}
			if len(userConfig.Pattern) > 0 {
				pattern = userConfig.Pattern
			}
		}
	}

	cFunc := icons.ColorV2
	if colorScheme == "v1" {
		cFunc = icons.ColorV1
	} else if colorScheme == "gh" || colorScheme == "github" {
		cFunc = icons.ColorGh
	}

	var iconGenerator icons.IconGenerator
	if pattern == "github" || pattern == "gh" {
		iconGenerator = icons.NewGhIconGenerator().WithColorGenerator(cFunc)
	} else {
		iconGenerator = icons.NewIdIconGenerator().WithColorGenerator(cFunc)
	}

	ct := r.URL.Query().Get("ct")
	cth := r.Header.Get("Accept")
	if ct == "svg" || cth == "image/svg+xml" {
		w.Header().Add("Content-Type", "image/svg+xml")
		_, err = w.Write([]byte(iconGenerator.GenSvg(id, size)))
	} else if ct == "jpeg" || cth == "image/jpeg" {
		w.Header().Add("Content-Type", "image/jpeg")
		err = jpeg.Encode(w, iconGenerator.GenIcon(id, size), nil)
	} else if ct == "gif" || cth == "image/gif" {
		w.Header().Add("Content-Type", "image/gif")
		err = gif.Encode(w, iconGenerator.GenIcon(id, size), nil)
	} else {
		w.Header().Add("Content-Type", "image/png")
		err = png.Encode(w, iconGenerator.GenIcon(id, size))
	}
}

var (
	config Config
)

func configure(configFile string) {
	if file, err := os.OpenFile(configFile, os.O_RDONLY, 0); err == nil {
		c := &Config{}
		_, err := toml.NewDecoder(file).Decode(c)
		if err == nil {
			log.Printf("Using configuration file '%s'", configFile)
		} else {
			log.Printf("Invalid configuration file '%s' - ignore it.\n", configFile)
		}

		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				log.Printf("Cannot close config file")
			}
		}(file)
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
	configFile := flag.String("c", "/etc/idicon/config.toml", "path to config file")
	port := flag.Int("p", 8000, "server port")

	flag.Parse()

	configure(*configFile)

	router := mux.NewRouter()
	router.HandleFunc("/avatar", pageRequestHandler)
	router.HandleFunc("/avatar/{id}", requestHandler)
	log.Printf("Starting on port %d ...\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), router))
}
