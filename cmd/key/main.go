package main

import (
	"fmt"
	"io"
	"os"

	"github.com/ryantate13/key_cli"
)

const (
	bin             = "key"
	defaultKeychain = "KEY_CLI"
	envVar          = defaultKeychain + "_KEYCHAIN"
)

// commands
type command int

const (
	get command = iota
	set
	rm
	ls
	env
)

func (c command) hasKey() bool {
	return c == get || c == set || c == rm
}

type args struct {
	command  command
	keychain string
	key      string
	file     string
	version  bool
	help     bool
}

func parseArgs(argv []string) args {
	var a args
	for i, arg := range argv {
		switch arg {
		case "-h", "--help":
			return args{help: true}
		case "-v", "--version":
			return args{version: true}
		case "-f", "--file":
			if i+1 < len(argv) {
				a.file = argv[i+1]
				argv[i], argv[i+1] = "", ""
			} else {
				return args{help: true}
			}
		}
	}
	argv = key_cli.Some(argv)
	if len(argv) == 0 {
		return args{help: true}
	}
	err := (&a.command).UnmarshalText([]byte(argv[0]))
	if err != nil {
		return args{help: true}
	}
	argv = argv[1:]
	k := key_cli.Coalesce(os.Getenv(envVar), defaultKeychain)
	if a.command.hasKey() {
		switch len(argv) {
		case 2:
			a.keychain = argv[0]
			a.key = argv[1]
		case 1:
			a.key = argv[0]
			a.keychain = k
		default:
			return args{help: true}
		}
	} else {
		switch len(argv) {
		case 1:
			a.keychain = argv[0]
		case 0:
			a.keychain = k
		default:
			return args{help: true}
		}
	}
	return a
}

func help() {
	fmt.Printf(`Usage: %s [options] [command] [args]

Get or set keys in the system keychain. The default keychain used is '%s'. This can be overridden with
the environment variable '%s' or passed as an extra argument to commands. For get and set, the
default is to print the value to stdout or read the value from stdin, respectively. This can be overridden
with the --file option.

Commands:
  get [keychain?] [key]  Get a key's value from the keychain
  set [keychain?] [key]  Set a key's value in the keychain
  rm  [keychain?] [key]  Remove a key/value from the keychain
  ls  [keychain?]        List all keys in the keychain
  env [keychain?]        Print keychain keys as environment variables. Example: 'export $(key env)'
    
Options:
  -f, --file     Path of file to read value from or write value to
  -v, --version  Show version number and quit
  -h, --help     Show this help message and quit
`, bin, defaultKeychain, envVar)
}

func version() {
	fmt.Printf("%s %s\n", bin, key_cli.Version)
}

func main() {
	opts := parseArgs(os.Args[1:])
	if opts.help {
		help()
		return
	}
	if opts.version {
		version()
		return
	}
	kc := key_cli.Must(key_cli.Open, opts.keychain)
	file := func(fallback *os.File) *os.File {
		if opts.file == "" || opts.file == "-" {
			return fallback
		}
		return key_cli.Must(os.Open, opts.file)
	}
	switch opts.command {
	case get:
		key_cli.Must(file(os.Stdout).Write, kc.Get(opts.key))
	case set:
		kc.Set(opts.key, key_cli.Must[io.Reader](io.ReadAll, file(os.Stdin)))
	case rm:
		kc.Remove(opts.key)
	case ls:
		for _, k := range kc.Keys() {
			fmt.Println(k)
		}
	case env:
		for _, k := range kc.Keys() {
			fmt.Printf("%s=%s\n", k, kc.Get(k))
		}
	}
}
