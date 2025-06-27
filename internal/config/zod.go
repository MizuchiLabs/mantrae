package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hypersequent/zen"
	"github.com/traefik/traefik/v3/pkg/config/dynamic"
)

// StructToZodSchema converts a struct to a zod schema (for use in the frontend)
func StructToZodSchema() {
	types := map[string]any{
		// Routers
		"httpRouter": dynamic.Router{},
		"tcpRouter":  dynamic.TCPRouter{},
		"udpRouter":  dynamic.UDPRouter{},

		// Services
		"httpService": dynamic.Service{},
		"tcpService":  dynamic.TCPService{},
		"udpService":  dynamic.UDPService{},

		// HTTP Middlewares
		"httpMiddleware": dynamic.Middleware{},

		// TCP Middlewares
		"tcpMiddleware": dynamic.TCPMiddleware{},
	}

	var builder strings.Builder

	// Add a header
	builder.WriteString("// This file is auto-generated via `zen.StructToZodSchema`.\n")
	builder.WriteString("// Do not edit manually.\n\n")
	builder.WriteString("import { z } from 'zod';\n\n")

	for _, strct := range types {
		schema := zen.StructToZodSchema(strct)
		builder.WriteString(fmt.Sprintf("%s\n", schema))
	}

	out := "./web/src/lib/gen/zen/traefik-schemas.ts"

	if err := os.MkdirAll(filepath.Dir(out), 0755); err != nil {
		panic(err)
	}
	if err := os.WriteFile(out, []byte(builder.String()), 0644); err != nil {
		panic(err)
	}

	fmt.Printf("✅ Zod schemas written to %s\n", out)
}
