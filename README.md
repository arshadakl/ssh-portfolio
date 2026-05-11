# ssh-portfolio

Interactive terminal portfolio built with Go + Bubbletea + Lipgloss.

## Build & Run

**Clone**
```bash
git clone https://github.com/arshadakl/ssh-portfolio.git
cd ssh-portfolio
```

**Run directly**
```bash
go run main.go
```

**Build**

```bash
# Linux / macOS
go build -o portfolio .
./portfolio

# Windows
go build -o portfolio.exe .
.\portfolio.exe
```

**Cross-compile**

```bash
# Linux (from any OS)
GOOS=linux GOARCH=amd64 go build -o portfolio .

# macOS (from any OS)
GOOS=darwin GOARCH=amd64 go build -o portfolio .

# Windows (from any OS)
GOOS=windows GOARCH=amd64 go build -o portfolio.exe .
```

## Requirements

Go 1.21+
