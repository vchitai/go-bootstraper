package configs

import (
	"embed"
)

type App struct {
	AssetFS     embed.FS
	AssetPath   string
	Version     string
}
