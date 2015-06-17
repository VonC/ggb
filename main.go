package main

import (
	"fmt"
	"os"

	"github.com/VonC/ggb/cmd"
)

func main() {
	fmt.Printf("ggb: ")
	err := cmd.RunCommand(os.Args)
	if err != nil {
		fmt.Printf("%s", err.Error())
		os.Exit(1)
	}
}

func build(args []string) error {
	fmt.Printf("build to be done with args '%v'", args)
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
