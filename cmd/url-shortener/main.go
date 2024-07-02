package main

import (
	"fmt"

	"github.com/MaKYaro/url-shortener/internal/config"
)

func main() {
	// TODO: parse config
	cfg := config.MustLoad()
	fmt.Println(cfg)

	// TODO: init logger

	// TODO: init storage

	// TODO: init router

	// TODO: run app
}
