package templates

import "embed"

//go:embed *.html
var TemplateFS embed.FS
