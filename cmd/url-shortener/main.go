package main

import (
	"fmt"

	"github.com/noctusha/url-shortener/internal/config"
)

func main() {
	// init config: cleanenv
	cfg := config.New()

	fmt.Println(cfg)

	// init logger: slog

	// init storage: postgres

	// init router: chi

	// run server
}
