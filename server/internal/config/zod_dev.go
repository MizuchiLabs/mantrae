//go:build dev
// +build dev

package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/hypersequent/zen"
	"github.com/traefik/traefik/v3/pkg/config/dynamic"
)

// StructToZodSchema converts a struct to a zod schema (for use in the frontend)
func StructToZodSchema() {
	types := map[string]any{
		"config": dynamic.Configuration{},
	}

	var builder strings.Builder

	// Add a header
	builder.WriteString("// This file is auto-generated via `zen.StructToZodSchema`.\n")
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

	out := "./web/src/lib/gen/zen/traefik-schemas.ts"

	if err := os.MkdirAll(filepath.Dir(out), 0755); err != nil {
		panic(err)
	}
	if err := os.WriteFile(out, []byte(builder.String()), 0644); err != nil {
		panic(err)
	}

	fmt.Printf("Zod schemas written to %s\n", out)
	os.Exit(0)
}
