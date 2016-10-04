// Package cmdmain provides simple subcommand support.
package cmdmain

import (
	"flag"
	"fmt"
	"os"
	"sort"
)

// Globals represent global flags that may be specified
// before subcommands on the command line.
var Globals *flag.FlagSet

func init() {
	Globals = flag.NewFlagSet("globals", flag.ContinueOnError)
	Globals.Usage = func() {}

	flag.CommandLine = Globals
}

type cmd struct {
	flags *flag.FlagSet
	cmd   Command
}

var cmds = make(map[string]cmd)

// Command represents regirsters a new command.
type Command interface {
	Run(args []string) error
	ArgNames() string
}

// Register regirsers a new command.
// The func makeCmd should initialize the flags needed for the Command,
// and return the new Command.
// It panics if name has already been registered.
func Register(name string, makeCmd func(flags *flag.FlagSet) Command) {
	if _, ok := cmds[name]; ok {
		panic(fmt.Errorf("command %q already registered", name))
	}

	flags := flag.NewFlagSet(os.Args[0]+" "+name, flag.ContinueOnError)
	flags.Usage = func() {}

	c := makeCmd(flags)
	if c == nil {
		panic(fmt.Errorf("Command for %q is nil", name))
	}
	cmds[name] = cmd{flags, c}
}

// Main parses arguments and runs registered Commands.
// It panics if no commands has been registered.
func Main() {
	if len(cmds) == 0 {
		panic("program has no commands defined")
	}

	err := Globals.Parse(os.Args[1:])
	if err != nil {
		if err != flag.ErrHelp {
			fmt.Fprintln(os.Stderr, err)
		}
		Usage()
	}

	args := Globals.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "No command specified.\n")
		Usage()
	}

	cmdname := args[0]
	c, ok := cmds[cmdname]
	if !ok {
		fmt.Fprintln(os.Stderr, "Unknown command: ", cmdname)
		Usage()
	}
	err = c.flags.Parse(args[1:])
	if err != nil {
		var cmdopts string
		if hasFlags(c.flags) {
			cmdopts = "[cmdopts] "
		}
		fmt.Fprintf(os.Stderr, "Usage:\n\n  %s [globalopts] %s %s%s\n\n",
			os.Args[0], cmdname, cmdopts, c.cmd.ArgNames())
		if hasFlags(c.flags) {
			fmt.Fprintf(os.Stderr, "%s options:\n", cmdname)
			c.flags.PrintDefaults()
		}
		os.Exit(2)
	}
	if err := c.cmd.Run(c.flags.Args()); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

// describer may be implemented by a Command
// and may be used to tell what it does.
type describer interface {
	Describe() string
}

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage:\n\n  %s [globalopts] cmd [cmdopts] [cmdargs]\n\n", os.Args[0])

	fmt.Fprintln(os.Stderr, "Commands:")
	var names []string
	for n := range cmds {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		var desc string
		if descr, ok := cmds[n].cmd.(describer); ok {
			desc = descr.Describe()
		} else {
			desc = "No description available."
		}
		fmt.Fprintf(os.Stderr, "  %s: %s\n", n, desc)
	}

	if hasFlags(Globals) {
		fmt.Fprintln(os.Stderr, "\nGlobal options:")
		Globals.PrintDefaults()
	}
	os.Exit(2)
}

func hasFlags(flags *flag.FlagSet) bool {
	any := false
	flags.VisitAll(func(*flag.Flag) {
		any = true
	})
	return any
}
