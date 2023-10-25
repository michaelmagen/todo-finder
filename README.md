# todo-finder

The Todo Finder CLI is a tool for locating and displaying todo comments in your codebase. It is built with Go and utilizes Cobra, Viper, and Bubble Tea to provide a seamless command-line interface.

## Installation

1. **Clone the Repository**:

```bash
   git clone https://github.com/michaelmagen/todo-finder
```

2. **Build the Application**:

```bash
cd todo-finder
go build -o todo-finder
```

3. **(Optional) Move the Binary to a Directory in Your PATH**:

```bash
mv todo-finder /usr/local/bin/
```

Now you should be able to run todo-finder from anywhere in your terminal.

## Usage

The CLI contains the following commands:

### List

```bash
todo-finder list <directory> [flags]
```

List all todo comments found in the directory. A specific directory can be passed in to search. If no directory is passed in, then the current directory is searched.

Available flags:

- `-a, --hidden`: Include hidden files/directories in the search.
- `-g, --no-gitignore`: Include files ignored by git.

### Marker

```bash
todo-finder marker <new marker>
```

Set the marker for a todo comment in the configuration file. The default marker for todo comments is "TODO:". Pass in a string that will be the new marker. If no string is passed in, it will print the current marker.

For more information on each command and its usage, run `todo-finder [command] --help`.
