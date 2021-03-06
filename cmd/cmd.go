package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	flag "github.com/spf13/pflag"
)

var commands = map[string]*Command{}

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

	// FlagSet for that command only
	fs *flag.FlagSet
	// FlagSet for that command and all sub-commands
	gfs *flag.FlagSet
	// FlagSet combined (command only and all sub-commands)
	afs *flag.FlagSet
	// output buffer for afs
	abuf *bytes.Buffer

	// function for adding FlagSet for that command
	ffs func(*flag.FlagSet)

	// function for adding flags for that command and any sub-Command FlagSet
	fgfs func(*flag.FlagSet)

	// Parent Command
	parent parent

	// Subcommands
	subcmds map[string]*Command

	// Args passed when running command
	args []string
}

type parent interface {
	fullCommand() string
	parseFlags()
	addSubCmd(*Command)
}

func NewCommand(name, usageLine, short, long string, run func([]string) error, parent *Command) *Command {
	cmd := &Command{
		name:      name,
		usageLine: usageLine,
		short:     short,
		long:      long,
		run:       run,
		fs:        flag.NewFlagSet(name, flag.ExitOnError),
		gfs:       flag.NewFlagSet(name, flag.ExitOnError),
		afs:       flag.NewFlagSet(name, flag.ContinueOnError),
		subcmds:   make(map[string]*Command),
		args:      []string{},
	}
	cmd.abuf = new(bytes.Buffer)
	cmd.afs.SetOutput(cmd.abuf)
	cmd.afs.Usage = cmd.fUsage
	// fmt.Printf("New Command '%s', nil parent: %v\n", cmd.name, cmd.parent == nil)
	if parent == nil {
		// This is a root command
		commands[name] = cmd
	} else {
		cmd.parent = *parent
		parent.addSubCmd(cmd)
	}
	// fmt.Printf("registered parent for name '%s': %v\n", name, commands)
	return cmd
}

func (cmd Command) addSubCmd(c *Command) {
	cmd.subcmds[c.name] = c
	cmd.fgfs(c.afs)
	cmd.fgfs(c.gfs)
}

// Set function for adding flags for that command and any sub-Command FlagSet
func (cmd *Command) SetGFS(fgfs func(*flag.FlagSet)) {
	cmd.fgfs = fgfs
	fgfs(cmd.afs)
	fgfs(cmd.gfs)
}

// Set function for adding flags for that command FlagSet only
func (cmd *Command) SetFS(ffs func(*flag.FlagSet)) {
	cmd.ffs = ffs
	ffs(cmd.afs)
	ffs(cmd.gfs)
}

// Runnable indicates this is a command that can be involved.
// Non runnable commands are only informational.
func (c *Command) Runnable() bool { return c.run != nil }

// RunCommand parses flags and runs the Command.
func RunCommand(args []string) error {
	// fmt.Printf("args %+v, c %v '%v'\n", os.Args, commands == nil, commands)
	cmd, err := commandFromArgs(os.Args)
	// fmt.Printf("cmd nil?? %v, args %+v, err='%v'\n", cmd == nil, os.Args, err)
	if err != nil {
		return err
	}
	cmd.parseFlags()
	if cmd.run == nil {
		return nil
	}
	return cmd.run(cmd.afs.Args())
}

func (cmd Command) parseFlags() {
	if cmd.parent != nil {
		cmd.parent.parseFlags()
	}
	// fmt.Printf("cmd name='%s': args='%+v' (afs='%+v')\n", cmd.name, cmd.args, cmd.afs)
	if err := cmd.afs.Parse(cmd.args); err != nil {
		fmt.Printf("Incorrect usage of %s:\n", cmd.fullCommand())
		fmt.Printf("%s", cmd.abuf.String())
		os.Exit(1)
	}
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
	// fmt.Printf("len %v\n", len(args))
	for i, arg := range args {
		if i == 0 {
			arg = filepath.Base(arg)
			ext := filepath.Ext(arg)
			if ext != "" {
				arg = arg[:len(arg)-len(ext)]
			}
			if strings.HasSuffix(arg, ".test") {
				return &Command{name: "test", afs: flag.NewFlagSet("test", flag.ContinueOnError)}, nil
			}
		}
		// fmt.Printf("arg='%s'\n", arg)
		if arg == "--" {
			cmd.args = append(cmd.args, args[i:]...)
			return cmd, nil
		}
		if strings.HasPrefix(arg, "-") {
			cmd.args = append(cmd.args, arg)
			continue
		}
		var subcmd *Command
		// fmt.Printf("cmd nil %v, reg %v\n", cmd == nil, commands)
		if cmd == nil {
			cmd = commands[arg]
			if cmd == nil {
				return nil, fmt.Errorf("Unknown command '%s'", arg)
			}
		} else {
			subcmd = cmd.subcmds[arg]
			if subcmd == nil {
				cmd.args = append(cmd.args, arg)
			} else {
				cmd = subcmd
			}
		}
	}
	// fmt.Printf("cmd found: '%+v'\n", cmd)
	if cmd == nil {
		return nil, fmt.Errorf("Unknown command from args '%v'", args)
	}
	return cmd, nil
}

func (c *Command) fUsage() {
	// fmt.Printf("fusage called for '%s'\n", c.name)
	s := strings.Split(c.abuf.String(), "\n")[0]
	c.abuf.Truncate(len(s))
	fmt.Fprintf(c.abuf, "\n%s", c.UsageFlags())
}

func (c *Command) UsageFlags() string {
	u := ""
	u = u + "local flags:\n"
	u = u + c.fs.FlagUsages()
	u = u + "global flags:\n"
	u = u + c.gfs.FlagUsages()
	return u
}

func (c *Command) Usage() string {
	// fmt.Printf("fusage called for '%s'\n", c.name)
	u := ""
	u = u + c.name + ": "
	u = u + c.short + "\n\n"
	u = u + c.usageLine + "\n"
	u = u + c.long + "\n\n"
	u = u + c.UsageFlags()
	return u
}
