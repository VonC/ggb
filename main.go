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

var cmdcurrent *cmd.Command

func build(args []string) error {
	cmdcurrent = cmdggb
	if verbose {
		fmt.Printf("build to be done with args '%v'", args)
	}
	checkGlobalFlag()
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
	_, gerr := prj.Golang("install -a .")
	// fmt.Printf("gout '%s', gerr '%v'\n", gout, gerr)
	if gerr != nil {
		fmt.Printf("%s", gerr.Error())
		os.Exit(1)
	}
	return nil
}

func add(args []string) error {
	cmdcurrent = cmdadd
	if verbose {
		fmt.Printf("Add to be done with args '%v'", args)
	}
	checkGlobalFlag()
	if len(args) == 0 {
		fmt.Println("At least one dependency url is expected")
		os.Exit(1)
	}
	for _, arg := range args {
		if err := addsub(arg); err != nil {
			return err
		}
	}
	return nil
}

func checkGlobalFlag() {
	if help {
		usage()
		os.Exit(0)
	}
	if debug {
		prj.Debug = true
	}
}

func usage() {
	fmt.Print(cmdadd.Usage())
}
