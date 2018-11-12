package cmdmain

type funcFlag struct {
	v bool
	f func()
}

var funcFlags []*funcFlag

// VersionFlag adds the version flag that calls showVersion
// when -version is specified on the command line.
func VersionFlag(showVersion func()) {
	FlagFunc("version", "show version", showVersion)
}

// FlagFunc adds a new function flag with the specified usage.
// If -flag is specified on the command line,
// Main calls f and exits.
func FlagFunc(flag, usage string, f func()) {
	ff := &funcFlag{
		f: f,
	}
	Globals.BoolVar(&ff.v, flag, false, usage)
	funcFlags = append(funcFlags, ff)
}
