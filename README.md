# sq

SPDX Query, A jq for SPDX 2.2 Querying Utility

Note for newbies: SPDX is a format where sw entities discloses what open source libraries they have used in building their software. Security expers and legal compliance experts uses this data to check if they have any license issues or security vulnerability is there..

https://en.m.wikipedia.org/wiki/Software_Package_Data_Exchange

## Description

A binary to query the spdx-sboms-JSON results.

By default uses _$HOME/sbom.spdx.json_ file to load the data. (you can pass custom \*.spdx.json file using _--config_ option any time)

## Sample

Display meta data with `sq meta` option
![Sample](https://github.com/dineshr93/sq/blob/main/screenshots/meta.png?raw=true)

Display pkgs list with `sq pkgs` option
![Sample](https://github.com/dineshr93/sq/blob/main/screenshots/sq_pkgs.png?raw=true)

Display files list with `sq files` option
![Sample](https://github.com/dineshr93/sq/blob/main/screenshots/files.png?raw=true)

Display spdx relationships list with `sq rels dig` option
![Sample](https://github.com/dineshr93/sq/blob/main/screenshots/dig.png?raw=true)

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
>sq rels dig -h
dig the relationship and show actual packages and files

For Ex:
To list all relationship: ./sq rels dig

Usage:
  sq rels dig [flags]

Flags:
  -h, --help   help for dig

Global Flags:
      --config string   YOUR SBOM JSON FILE (default is $CURRENT_DIR/sbom.spdx.json)
=================================================================
Alternatively if UI is small to fit every thing, you can save the output to the file

sq meta > sbom-meta.txt
sq files > sbom-files.txt
sq pkgs > sbom-pkgs.txt
sq rels > sbom-rels.txt
```

## Authors

Dinesh Ravi

## Version History

- 1.0.0
  - Initial Release

## License

This project is licensed under the Apache License 2.0 - see the [Apache-2.0](LICENSE) file for details

## Acknowledgments

- [cobra](https://www.github.com/spf13/cobra)
- [viper](https://www.github.com/spf13/viper)
- [simpletable](https://www.github.com/alexeyco/simpletable)
