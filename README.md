# sq

An SPDX Querying Binary

## Description

A binary to query the sbom results.

By default uses _$HOME/sbom.spdx.json_ file to load the data. (you can pass custom *.spdx.json file using _--config_ option any time)

## Sample

![Sample](https://github.com/dineshr93/sq/blob/main/sample.png?raw=true)

## Getting Started

Contains following commands

1. List Files (todo add)
2. List Relationships (todo ls)
3. Delete (supports space separated multiple params) (todo delete)
4. Mark Complete (supports space separated multiple params) (todo complete)
5. Mark Pending (supports space separated multiple params) (todo pending)

### Dependencies

- Cobra
- Viper
- Simple table

### Installing

Choose appropriate (binary Releases)[https://github.com/dineshr93/todo/releases]
- Rename the binary to 'todo'.
- Add the binary to your environment path and use it. 

### Executing program

- How to run the program

```
COMMAND FORMAT:
todo add "project - task - jiraLink"
todo add "YOUR_PROJECT_NAME - TASK_NAME - http://www.yourcompany.jira.com"

todo complete TASKID
<!-- you can also mark multiple tasks to mark complete by passing space separated ids  -->
todo complete 1 2 ..

todo pending TASKID
<!-- you can also mark multiple tasks back to mark pending by passing space separated ids  -->
todo pending 1 2 ..

todo delete TASKID
<!-- you can also mark multiple tasks for deletion by passing space separated ids  -->
todo delete 1 2 ..
```

## Help

any issues please raise an bug reguest

```
>todo -h

        TODO CLI APP (for issue -> https://github.com/dineshr93/todo/issues)

1. Add (todo add)
2. List (todo ls)
3. Delete (supports space separated multiple params) (todo delete)
4. Mark Complete (supports space separated multiple params) (todo complete)
5. Mark Pending (supports space separated multiple params) (todo pending)

Usage:
  todo [command]

Available Commands:
  add         add a todo Format: Project - task - JIRA link
  complete    Mark a task as completed(supports space separated multiple params)
  delete      Delete a task from the list(supports space separated multiple params)
  help        Help about any command
  ls          List all todos
  pending     Mark a task as pending(supports space separated multiple params)

Flags:
      --config string   config file (default is $HOME/.todo.json)
  -h, --help            help for todo
  -t, --toggle          Help message for toggle

Use "todo [command] --help" for more information about a command.

```

## Authors

Dinesh Ravi

## Version History

- 1.0.0
  - Initial Release

## License

This project is licensed under the Apache License 2.0 or GPL-3.0 - see the [Apache-2.0](LICENSE) file for details

## Acknowledgments

- [cobra](https://www.github.com/spf13/cobra)
- [viper](https://www.github.com/spf13/viper)
- [simpletable](https://www.github.com/alexeyco/simpletable)
- [joefazee](https://www.github.com/joefazee) (for models under GPL-2.0)
