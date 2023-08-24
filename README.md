# Quake Log Stats 

## Description
This program consists of a CLI that parses a Quake 3 Arena log file and generates a report with the following information for each game:
- Total kills
- Players
- Kills by means
- Kills by players

## How to run

### Go run
To run the CLI just run the command `go run . file <log_file_name>` in the root of the project.

Example:
```bash
go run . -file games.log
```

### Go build
Run the command `go build .` in the root of the project and then run the generated binary.