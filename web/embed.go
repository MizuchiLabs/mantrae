package web

import (
	"embed"
)

//go:embed all:build
var StaticFS embed.FS
