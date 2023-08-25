# Quake Log Stats 

![Preview](preview.png)

## Description
This program consists of a CLI that parses a Quake 3 Arena log file and generates a data file with the following information for each game:
- Total kills
- Players
- Kills by means
- Kills by players

and a report in HTML format with the following information for each game:
- Total kills
- Players
- Kills by means
- Kills by players
- Players ranking based on kills

You can find samples of the generated files in the `samples` folder.

The report will also be printed in the terminal on the following format:
```go
------------------------------------------------
Match 20
Total kills: 3
Players:
        Oootsimo: 1
        Dono da Bola: 2
        Assasinu Credi: 0
        Zeh: 0
        Mal: 0
        Isgalamido: 0
Kills by means of death:
        MOD_ROCKET_SPLASH: 2
        MOD_ROCKET: 1
Ranking:
        1. Dono da Bola: 2
        2. Oootsimo: 1
        3. Mal: 0
        4. Isgalamido: 0
        5. Assasinu Credi: 0
        6. Zeh: 0

```
## How to run

### Go run
To run the CLI just run the command `go run . file <log_file_name>` in the root of the project.

Example:
```bash
go run . -file games.log
```

### Go build
Run the command `go build .` in the root of the project and then run the generated binary.