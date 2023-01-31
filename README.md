# sq

SPDX Query, A jq for SPDX 2.2 Querying Utility

## Description

A binary to query the spdx-sboms-JSON results.

By default uses _$HOME/sbom.spdx.json_ file to load the data. (you can pass custom \*.spdx.json file using _--config_ option any time)

## Sample

![Sample](https://github.com/dineshr93/sq/blob/main/sample.png?raw=true)

## Getting Started

Contains following commands

1. List Meta ata (sq meta)
2. List Relationships (sq ls)
3. Delete (supports space separated multiple params) (sq delete)
4. Mark Complete (supports space separated multiple params) (sq complete)
5. Mark Pending (supports space separated multiple params) (sq pending)

### Dependencies

- Cobra
- Viper
- Simple table

### Installing

Choose appropriate (binary Releases)[https://github.com/dineshr93/sq/releases]

- Rename the binary to 'sq'.
- Add the binary to your environment path and use it.

### Executing program

- How to run the program

```
>sq -h
A SBOM Query CLI (for issue -> https://github.com/dineshr93/sq/issues)

        1. List Meta data (sq meta)
        2. List Files
        3. List Packages
        4. List Relationships

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
