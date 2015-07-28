package main

import (
	flag "github.com/spf13/pflag"

	"github.com/VonC/ggb/cmd"
)

var debug bool
var verbose bool
var help bool
var version bool
var cmdggb *cmd.Command
var cmdadd *cmd.Command

func init() {
	initCommands()
}

func initCommands() {
	initCommandGgb()
	initCommandAdd()
}

func initCommandGgb() {
	cmdggb = cmd.NewCommand(
		"ggb",
		"ggb [cmd]",
		"builds a go project with git submodule dependencies management",
		`ggb (alone) builds, 
while 'ggb deps' offers dependency management as git submodules`,
		build, nil)
	cmdggb.SetGFS(func(gfs *flag.FlagSet) {
		gfs.BoolVarP(&help, "help", "h", false, "ggb usage")
		gfs.BoolVarP(&verbose, "verbose", "v", false, `display a verbose output
		not suited for batch usage`)
		gfs.BoolVarP(&debug, "debug", "d", false, "output debug informations (not for batch usage)")
		gfs.BoolVarP(&version, "version", "V", false, "display ggb version")
	})
}

func initCommandAdd() {
	cmdadd = cmd.NewCommand(
		"add",
		"ggb add [cmd]",
		"Add a dependency to a go project with git submodule dependencies management",
		`ggb add will add as a submodule a dependency,
 and check the dependencies of the dependency added`,
		add, cmdggb)
}
