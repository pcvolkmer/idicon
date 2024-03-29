package main

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/gorilla/mux"
	"idicon/icons"
	"image/png"
	"log"
	"net/http"
	"os"
	"strconv"
)

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

	w.Header().Add("Content-Type", "image/png")
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

	err = png.Encode(w, iconGenerator.GenIcon(id, size))
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
	configFile := flag.String("c", "/etc/idicon/config.toml", "path to config file")
	port := flag.Int("p", 8000, "server port")

	flag.Parse()

	configure(*configFile)

	router := mux.NewRouter()
	router.HandleFunc("/avatar/{id}", requestHandler)
	log.Printf("Starting on port %d ...\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), router))
}
