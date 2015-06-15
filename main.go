package main

import (
	"fmt"
	"os"
)
import flag "github.com/spf13/pflag"

var help bool

func init() {
	flag.BoolVarP(&help, "help", "h", false, "ghref usage")
}

func main() {
	flag.Parse()
	if help {
		usage()
		os.Exit(0)
	}
	fmt.Printf("ggb: ")
	if len(flag.Args()) == 0 {
		build()
	} else {
		processCommands(flag.Args())
	}
}

func processCommands(args []string) {
	switch args[0] {
	case "addsub":
		addsub(args)
	default:
		fmt.Print("Unknown command " + args[0])
	}
}

func build() {
	fmt.Print("build to be done")
}

func usage() {
	fmt.Print(`ggb [-h],
builds a go project with git submodule dependencies`)
}
