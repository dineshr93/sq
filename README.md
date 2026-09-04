# sq

SPDX Query - Making SPDX (2.2, 2.3 & 3.0.1) JSON Files human readable

Note for newbies: SPDX is a format where sw entities discloses what open source libraries they have used in building their software. Security expers and legal compliance experts uses this data to check if they have any license issues or security vulnerability is there..

https://en.m.wikipedia.org/wiki/Software_Package_Data_Exchange

## Docker image
```
docker pull dineshr93/sq:1.0
```

## Load alias

```
alias dr='docker run'
alias p='echo ${PWD}'
alias sq='dr -v ${PWD}:${PWD} dineshr93/sq:1.0'
```

## command

```
sq -c $(p)/ubuntu20.04.spdx.json -h
```

### example
```
sq -c $(p)/ubuntu20.04.spdx.json pkgs 5

```

## Description

A binary to query the spdx-sbom JSON results.

SPDX 3.0.1 documents are auto-detected on load and rendered through the same
commands; 3.0-only constructs that have no 2.x equivalent are counted in
`sq meta` under *3.0 elements w/o 2.x mapping*.

By default uses _$HOME/sbom.spdx.json_ file to load the data. (you can pass custom \*.spdx.json file using _--config_ option any time)

## Sample

Display meta data with `sq meta` option
![Sample](https://github.com/dineshr93/sq/blob/main/screenshots/meta.png?raw=true)

If `--config` option is not passed it will detect & load first spdx json file automatically
Display pkgs list with `sq pkgs` option
![Sample](https://github.com/dineshr93/sq/blob/main/screenshots/noconfig.png?raw=true)

limit pkgs with `sq pkgs NUMBER` option
![Sample](https://github.com/dineshr93/sq/blob/main/screenshots/sq_pkgs.png?raw=true)

Display files list with `sq files` option
![Sample](https://github.com/dineshr93/sq/blob/main/screenshots/files.png?raw=true)

Display spdx relationships table and list with `sq rels` option
![Sample](https://github.com/dineshr93/sq/blob/main/screenshots/rels.png?raw=true)

Display spdx relationships list with `sq rels dig` option
![Sample](https://github.com/dineshr93/sq/blob/main/screenshots/dig.png?raw=true)

Display IP Details list with `sq pkgs ip` option
![Sample](https://github.com/dineshr93/sq/blob/main/screenshots/ip.png?raw=true)

## Getting Started

Contains following commands

        1. List Meta data (sq meta)
        2. List Files(sq files)
        3. List Packages (sq pkgs)
        4. List Relationships (sq rels)
        5. List pkgs and files in Relationships`(sq rels dig)

### Dependencies

- Cobra
- Viper
- Simple table
- tools-golang (SPDX 3.0.1 parser)

### Installing

Choose appropriate (binary Releases)[https://github.com/dineshr93/sq/releases]

- Rename the binary to 'sq'.
- Add the binary to your environment path and use it.

### How to run

- How to run the program

```
>sq -h
A SBOM Query CLI (for issue -> https://github.com/dineshr93/sq/issues)

        1. List Meta data (sq meta)
        2. List Files (sq files)
        3. List Packages (sq pkgs)
        4. List Relationships (sq rels)
        5. List pkgs and files in Relationships (sq rels dig)

Usage:
  sq [command]

Available Commands:
  files       Command to list files section
  help        Help about any command
  meta        Meta data of the spdx file
  pkgs        Command to list pkgs section
  rels        Lists Relationships

Flags:
      --config string   config file (default is $HOME/.sq.yaml)
  -h, --help            help for sq
  -t, --toggle          Help message for toggle

Use "sq [command] --help" for more information about a command.
================================================================
Alternatively if UI is small to fit every thing, you can save the output to the file

sq meta > sbom-meta.txt
sq files > sbom-files.txt
sq pkgs > sbom-pkgs.txt
sq rels > sbom-rels.txt
```

## Build from source

Requires Go 1.23.5+

```
git clone https://github.com/dineshr93/sq && cd sq
make build    # builds ./bin/sq and copies it to ./sq
go test ./... # run the test suite
```

## Sample data

Grab real SPDX 3.0.1 SBOMs from the official examples repo:

```
mkdir -p /tmp/sq-samples && cd /tmp/sq-samples
B=raw.githubusercontent.com/spdx/spdx-examples/master
curl -fsSL -o example1.spdx3.json                    $B/software/example1/spdx3.0/example1.spdx3.json
curl -fsSL -o example3-bin.spdx3.json                $B/software/example3/spdx3.0/example3-bin.spdx3.json
curl -fsSL -o examplemaven-0.0.1-enriched.spdx3.json $B/software/example14/spdx3.0/examplemaven-0.0.1-enriched.spdx3.json
```

Then point sq at any of them:

```
./sq meta --config /tmp/sq-samples/example1.spdx3.json
./sq pkgs --config /tmp/sq-samples/examplemaven-0.0.1-enriched.spdx3.json
```

Or drop the files into your current folder - sq auto-detects and loads the
first valid SPDX JSON without `--config`.

## Authors

Dinesh Ravi

## Version History

- 1.1.0
  - SPDX 3.0.1 support (auto-detected; all commands work on 2.x and 3.0.1)
- 1.0.0
  - Initial Release

## License

This project is licensed under the Apache License 2.0 - see the [Apache-2.0](LICENSE) file for details

## Acknowledgments

- [cobra](https://www.github.com/spf13/cobra)
- [viper](https://www.github.com/spf13/viper)
- [simpletable](https://www.github.com/alexeyco/simpletable)
