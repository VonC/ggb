package main

import (
	"fmt"
	"os"

	"github.com/VonC/ggb/cmd"
	"github.com/VonC/ggb/prj"
)

func main() {
	if verbose {
		fmt.Printf("ggb: ")
	}
	err := cmd.RunCommand(os.Args)
	if err != nil {
		fmt.Printf("%s", err.Error())
		os.Exit(1)
	}
}

func build(args []string) error {
	if verbose {
		fmt.Printf("build to be done with args '%v'", args)
	}
	p, err := prj.GetProject()
	if err != nil {
		fmt.Printf("%s", err.Error())
		os.Exit(1)
	}
	if verbose {
		fmt.Printf(" in root folder '%s'\n", p.RootFolder())
	}
	// gout, gerr := prj.Golang("env")
	// fmt.Printf("gout '%s', gerr '%v'\n", gout, gerr)
	prj.Golang("install -a .")
	// fmt.Printf("gout '%s', gerr '%v'\n", gout, gerr)
	return nil
}

func checkGlobalFlag() {
	if help {
		usage()
		os.Exit(0)
	}
}

func usage() {
	fmt.Print(`ggb [-h],
builds a go project with git submodule dependencies`)
}
