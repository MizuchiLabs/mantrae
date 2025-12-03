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

	withDesc := map[string]zen.CustomFn{
		"zodDesc": func(c *zen.Converter, t reflect.Type, desc string, indent int) string {
			return fmt.Sprintf(".describe(%q)", desc)
		},
	}

	for _, strct := range types {
		schema := zen.StructToZodSchema(strct, zen.WithCustomTags(withDesc))
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
