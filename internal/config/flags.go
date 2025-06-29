// Package config provides functions for parsing command-line f and
// setting up the application's default settings.
package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/mizuchilabs/mantrae/internal/store"
	"github.com/mizuchilabs/mantrae/pkg/build"
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
	// flag.BoolVar(&f.Squash, "squash", false, "Squash the database (only for dev)")
	// flag.BoolVar(&f.Zod, "zod", false, "Generate zod schemas (only for dev)")

	flag.Parse()
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	if f.Version {
		fmt.Println(build.Version)
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
