# `key-cli`

## Usage:

```console
key [options] [command] [args]
```

Get or set keys in the system keychain. The default keychain used is `KEY_CLI`. This can be overridden with
the environment variable `KEY_CLI_KEYCHAIN` or passed as an extra argument to commands. For `get` and `set`, the
default is to print the value to stdout or read the value from stdin, respectively. This can be overridden
with the `--file` option.

### Commands
 - `get` [keychain?] [key]  Get a key`s value from the keychain
 - `set` [keychain?] [key]  Set a key`s value in the keychain
 - `rm`  [keychain?] [key]  Remove a key/value from the keychain
 - `ls`  [keychain?]        List all keys in the keychain
 - `env` [keychain?]        Print keychain keys as environment variables. Example: `export $(key env)`
    
### Options
 - `-f`, `--file`     Path of file to read value from or write value to
 - `-v`, `--version`  Show version number and quit
 - `-h`, `--help`     Show this help message and quit
