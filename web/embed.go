package web

import (
	"embed"
)

//go:generate pnpm install
//go:generate pnpm build
//go:embed all:build
var StaticFS embed.FS
