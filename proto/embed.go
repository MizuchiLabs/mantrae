package proto

import (
	"embed"
)

//go:embed gen/openapi/openapi.yaml
var OpenAPIFS embed.FS
