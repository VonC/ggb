package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	flag "github.com/spf13/pflag"
)

var commands map[string]*Command = make(map[string]*Command)

// inspired by https://github.com/constabulary/gb/blob/master/cmd/cmd.go

// Command represents a subcommand, or plugin that is executed
type Command struct {
	// Name of the command
	name string

	// UsageLine demonstrates how to use this command
	usageLine string

	// Single line description of the purpose of the command
	short string

	// Description of this command
	long string

	// Run is invoked with arguments left over after flag parsing.
	run func(args []string) error

	// FlagSet for adding flags for that command
	fs *flag.FlagSet

	// FlagSet for adding flags for that command and any sub-Command
	gfs *flag.FlagSet

	// Parent Command
	parent parent

	// Subcommands
	subcmds map[string]*Command
}

type parent interface {
	fullCommand() string
	checkGlobalFlags(args []string) []string
}

func NewCommand(name, usageLine, short, long string, run func([]string) error, parent *Command) *Command {
	cmd := &Command{
		name:      name,
		usageLine: usageLine,
		short:     short,
		long:      long,
		run:       run,
		fs:        flag.NewFlagSet(name, flag.ExitOnError),
		gfs:       flag.NewFlagSet(name, flag.ContinueOnError),
		subcmds:   make(map[string]*Command),
	}
	if parent == nil {
		// This is a root command
		commands[name] = cmd
	} else {
		cmd.parent = *parent
	}
	return cmd
}

// FlagSet for adding flags for that command
func (cmd *Command) FS() *flag.FlagSet {
	return cmd.fs
}

// FlagSet for adding flags for that command and any sub-Command
func (cmd *Command) GFS() *flag.FlagSet {
	return cmd.gfs
}

// Runnable indicates this is a command that can be involved.
// Non runnable commands are only informational.
func (c *Command) Runnable() bool { return c.run != nil }

// RunCommand parses flags and runs the Command.
func RunCommand(args []string) error {
	cmd, err := commandFromArgs(os.Args)
	if err != nil {
		return err
	}
	args = cmd.checkGlobalFlags(args)
	fs := cmd.fs
	if err := fs.Parse(args); err != nil {
		fmt.Printf("Incorrect usage of %s:", cmd.fullCommand())
		fs.Usage()
		os.Exit(1)
	}
	if cmd.run == nil {
		return nil
	}
	return cmd.run(cmd.fs.Args()[1:])
}

func (cmd Command) fullCommand() string {
	res := cmd.name
	if cmd.parent != nil {
		res = cmd.parent.fullCommand() + " " + res
	}
	return res
}

func commandFromArgs(args []string) (*Command, error) {
	var cmd *Command
	for i, arg := range args {
		if i == 0 {
			arg = filepath.Base(arg)
			ext := filepath.Ext(arg)
			if ext != "" {
				arg = arg[:len(arg)-len(ext)]
			}
		}
		if arg == "--" {
			return cmd, nil
		}
		if strings.HasPrefix(arg, "-") {
			continue
		}
		var subcmd *Command
		if cmd == nil {
			cmd = commands[arg]
		} else {
			subcmd = cmd.subcmds[arg]
		}
		if cmd == nil {
			return nil, fmt.Errorf("Unknown command '%s'", arg)
		}
		if subcmd == nil {
			// let's consider all the rest as parameters of the current command
			return cmd, nil
		}
		cmd = subcmd
	}
	if cmd == nil {
		return nil, fmt.Errorf("Unknown command from args '%v'", args)
	}
	return cmd, nil
}

func (cmd Command) checkGlobalFlags(args []string) []string {
	if cmd.parent != nil {
		args = cmd.parent.checkGlobalFlags(args)
	}
	cmd.gfs.Parse(args)
	args = cmd.gfs.Args()
	return args
}
