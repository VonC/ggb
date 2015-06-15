package main

import "./cmd"

var commands map[string]*cmd.Command

func init() {
	initCommands()
}

func initCommands() {
	initCommandGgb()
}

func initCommandGgb() {
	cmd := cmd.NewCommand("ggb",
		"ggb [cmd]",
		"builds a go project with git submodule dependencies management",
		`ggb builds, 
while ggb deps offers dependency management as git submodules`,
		nil)
	commands["ggb"] = cmd
}
