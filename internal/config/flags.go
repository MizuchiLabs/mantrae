// Package config provides functions for parsing command-line f and
// setting up the application's default settings.
package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/MizuchiLabs/mantrae/internal/util"
)

type Flags struct {
	Version bool
	Update  bool
	Reset   bool
}

func ParseFlags() (*Flags, error) {
	f := &Flags{}
	flag.BoolVar(&f.Version, "version", false, "Print version and exit")
	flag.BoolVar(&f.Update, "update", false, "Update the application")
	flag.BoolVar(&f.Reset, "reset", false, "Reset the default admin password")

	flag.Parse()
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	if f.Version {
		fmt.Println(util.Version)
		os.Exit(0)
	}

	return f, nil
}
