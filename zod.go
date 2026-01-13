package main

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/hypersequent/zen"
	"github.com/traefik/traefik/v3/pkg/config/dynamic"
)

const (
	output = "./web/src/lib/gen/zen/traefik-schemas.ts"
)

func main() {
	types := map[string]any{"config": dynamic.Configuration{}}
	var builder strings.Builder

	// Add a header
	builder.WriteString("// This file is auto-generated via `zen`.\n")
	builder.WriteString("// Do not edit manually.\n\n")
	builder.WriteString("import { z } from 'zod';\n\n")

	customTypeHandlers := map[string]zen.CustomFn{
		"github.com/traefik/paerser/types.Duration": func(c *zen.Converter, t reflect.Type, v string, indent int) string {
			return "z.string()"
		},
	}

	for _, strct := range types {
		schema := zen.StructToZodSchema(strct, zen.WithCustomTypes(customTypeHandlers))
		builder.WriteString(fmt.Sprintf("%s\n", schema))
	}

	if err := os.MkdirAll(filepath.Dir(output), 0o750); err != nil {
		panic(err)
	}
	if err := os.WriteFile(output, []byte(builder.String()), 0o600); err != nil {
		panic(err)
	}

	fmt.Printf("Zod schemas written to %s\n", output)
	os.Exit(0)
}
