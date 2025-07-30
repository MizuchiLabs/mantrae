//go:build dev
// +build dev

package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/mizuchilabs/mantrae/pkg/build"
	"github.com/mizuchilabs/mantrae/pkg/meta"
	"github.com/mizuchilabs/mantrae/server/internal/store"
)

type Flags struct {
	Version bool
	Update  bool
	Squash  bool
	Zod     bool
}

func ParseFlags() {
	f := &Flags{}
	flag.BoolVar(&f.Version, "version", false, "Print version and exit")
	flag.BoolVar(&f.Update, "update", false, "Update the application")
	flag.BoolVar(&f.Squash, "squash", false, "Squash the database (only for dev)")
	flag.BoolVar(&f.Zod, "zod", false, "Generate zod schemas (only for dev)")

	flag.Parse()
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	if f.Version {
		fmt.Println(meta.Version)
		os.Exit(0)
	}

	if f.Squash {
		store.Squash()
		os.Exit(1)
	}

	if f.Zod {
		StructToZodSchema()
		os.Exit(1)
	}
	build.Update(f.Update)
}
