package main

import (
	"github.com/yaroslavst777/calculator/internal/application"
)

func main() {
	app := application.New()
	app.RunServer()
}
