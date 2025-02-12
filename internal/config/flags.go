// Package config provides functions for parsing command-line f and
// setting up the application's default settings.
package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/MizuchiLabs/mantrae/internal/db"
	"github.com/MizuchiLabs/mantrae/pkg/build"
)

type Flags struct {
	Version bool
	Update  bool
	Squash  bool
}

func ParseFlags() {
	f := &Flags{}
	flag.BoolVar(&f.Version, "version", false, "Print version and exit")
	flag.BoolVar(&f.Update, "update", false, "Update the application")
	flag.BoolVar(&f.Squash, "squash", false, "Squash the database")

	flag.Parse()
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	if f.Version {
		fmt.Println(build.Version)
		os.Exit(0)
	}

	if f.Squash {
		db.Squash()
		os.Exit(1)
	}
	build.Update(f.Update)
}
