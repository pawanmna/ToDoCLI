# ToDoCLI

A fast Go-based command-line task manager with SQLite persistence.

## Features
- Add tasks
- List pending/all tasks
- Mark tasks done
- Delete tasks

## Install
### From source
```bash
1. Build and run locally:

go build -o todo.exe .
.\todo.exe -add "task"

2. Install to Go bin:

go install .
Then run:
todo.exe -add "task"

This only works if your Go bin directory is in PATH.
```

## Usage
```bash
todo -add "Buy milk"
todo -list
todo -done -id 1
todo -delete -id 1
todo -list -all
```

## Example Output
<sample terminal output>

## Roadmap
- Tags
- Priorities
- Due dates
- Subtasks
- AI organization

