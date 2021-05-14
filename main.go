package main

import (
	"embed"
	"log"

	"go-bootstraper/configs"
	"go-bootstraper/internal/cli"
)

//go:embed assets/templates
var assets embed.FS
//go:embed assets/version
var version []byte

func main() {
	app, err := cli.NewApp(&configs.App{
		AssetFS:   assets,
		AssetPath: "assets/templates",
		Version:   string(version),
	})
	if err != nil {
		log.Fatal(err)
	}
	app.Run()
}
