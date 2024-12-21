package main

import (
	"github.com/yaroslavst777/Calculator/internal/application"
)

func main() {
	app := application.New()
	app.RunServer()
}
