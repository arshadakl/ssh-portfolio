# ssh-portfolio

Interactive terminal portfolio built with Go + Bubbletea + Lipgloss.

## Try it

```bash
ssh arshadakl.in
```

## Controls

| Key | Action |
| --- | --- |
| `j/k` or `↑/↓` | Navigate sections |
| `w/s` or mouse wheel | Scroll content |
| `1-9` | Jump to section |
| `:` | Open command bar |
| `q` | Quit |

In the **contribute** section: `←/→` or `tab` chooses a visitor intent, `enter` expands it, `esc` collapses.

In the **ask-ai** section, typing goes directly into the AI prompt. `enter` sends, `esc` returns to portfolio navigation, and `/help` lists the local AI commands.

### Commands

Open the command bar with `:`, then:

- **Sections** — `home` `ask-ai` `whoami` `experience` `projects` `contribute` `recognition` `skills` `contact`
- **Links** — `resume` `email` `github` `linkedin` `blog` `leetcode` (prints a copyable hyperlink)
- **Shortcuts** — `hire` (jumps to contribute with "hiring a developer" expanded), `clear`, `help`, `quit`
- **Contact form** — `message` or `msg` opens a step-by-step wizard that POSTs to the site's contact API

## Run locally

```bash
# generate a host key (once)
ssh-keygen -t ed25519 -f keys/id_ed25519

# start the server on port 2222
PORT=2222 HOST_KEY_PATH=keys/id_ed25519 go run .

# connect
ssh -p 2222 localhost
```

The AI section calls the deployed Portfolio RAG Worker. Configure the trusted client token before starting:

```bash
PORTFOLIO_RAG_URL=https://portfolio-rag.arshad-v2.workers.dev
PORTFOLIO_RAG_SSH_TOKEN=your-shared-secret
```

Both variables are required. The Go process reads them from its runtime environment; it does not contain a fallback Worker URL and does not automatically load `.env` files.

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
