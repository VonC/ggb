package cmd

import (
	"os"

	flag "github.com/spf13/pflag"
)

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

	// FlagSet for adding flags
	fs *flag.FlagSet
}

func NewCommand(name, usageLine, short, long string, run func([]string) error) *Command {
	cmd := &Command{
		name:      name,
		usageLine: usageLine,
		short:     short,
		long:      long,
		run:       run,
		fs:        flag.NewFlagSet(name, flag.ExitOnError),
	}
	return cmd
}

// FlagSet for adding flags
func (cmd *Command) FS() *flag.FlagSet {
	return cmd.fs
}

// Runnable indicates this is a command that can be involved.
// Non runnable commands are only informational.
func (c *Command) Runnable() bool { return c.run != nil }

// RunCommand parses flags and runs the Command.
func RunCommand(cmd *Command, args []string) error {
	fs := cmd.FS()
	cmdflags, args := splitargs(args)
	if err := fs.Parse(cmdflags); err != nil {
		fs.Usage()
		os.Exit(1)
	}
	// TODO check subcommand
	return cmd.run(args)
}

func splitargs(args []string) ([]string, []string) {
	return []string{}, []string{}
}
