package main

import "./cmd"

var debug bool
var verbose bool
var help bool
var version bool

func init() {
	initCommands()
}

func initCommands() {
	initCommandGgb()
}

func initCommandGgb() {
	cmdggb := cmd.NewCommand("ggb",
		"ggb [cmd]",
		"builds a go project with git submodule dependencies management",
		`ggb builds, 
while ggb deps offers dependency management as git submodules`,
		build, nil)

	gfs := cmdggb.GFS()
	gfs.BoolVarP(&help, "help", "h", false, "ggb usage")
	gfs.BoolVarP(&verbose, "verbose", "v", false, `display a verbose output
		not suited for batch usage`)
	gfs.BoolVarP(&debug, "debug", "d", false, "output debug informations (not for batch usage)")
	gfs.BoolVarP(&version, "version", "V", false, "display ggb version")
}
