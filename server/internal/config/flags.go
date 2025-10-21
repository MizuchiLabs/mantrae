//go:build !dev
// +build !dev

// Package config various app setup and configuration functions.
package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/mizuchilabs/mantrae/pkg/build"
	"github.com/mizuchilabs/mantrae/pkg/meta"
)

type Flags struct {
	Version       bool
	Update        bool
	Squash        bool
	Zod           bool
	ResetPassword string
	ResetUser     string
}

func ParseFlags() *Flags {
	f := &Flags{}
	flag.BoolVar(&f.Version, "version", false, "Print version and exit")
	flag.BoolVar(&f.Update, "update", false, "Update the application")
	flag.StringVar(&f.ResetPassword, "reset-password", "", "Set a new admin password")
	flag.StringVar(&f.ResetUser, "reset-user", "admin", "Username to reset password for")

	flag.Parse()
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	if f.Version {
		fmt.Println(meta.Version)
		os.Exit(0)
	}

	build.Update(f.Update)
	return f
}
