# LazyCommands

A terminal UI tool for running shell commands sequentially with visual progress tracking.

## Features

- **Simple single-panel UI**: Shows all commands with status indicators
- **Sequential execution**: Runs commands one after another
- **Status tracking**: Visual indicators for pending, running, completed, and failed commands
- **Automatic termination on failure**: Stops execution and displays error output if a command fails
- **Clean and fast**: Minimal interface focused on getting work done

## Prerequisites

### Installing Go

LazyCommands is built with Go. You need Go 1.21 or later installed on your system.

#### macOS

**Option 1: Using Homebrew (Recommended)**
```bash
brew install go
```

**Option 2: Direct Download**
1. Visit [https://go.dev/dl/](https://go.dev/dl/)
2. Download the macOS installer
3. Run the installer and follow the prompts

#### Linux

**Option 1: Using package manager**
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install golang-go

# Fedora
sudo dnf install golang

# Arch
sudo pacman -S go
```

**Option 2: Direct Download**
```bash
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
```

#### Windows

1. Visit [https://go.dev/dl/](https://go.dev/dl/)
2. Download the Windows installer (.msi file)
3. Run the installer

#### Verify Installation

```bash
go version
```

You should see output like: `go version go1.21.0 darwin/arm64`

## Installation

Once Go is installed, you can install LazyCommands:

### From Source

```bash
cd /path/to/lazycommands
go install
```

This will install the `lazycommands` binary to `$GOPATH/bin` (usually `~/go/bin`).

Make sure `$GOPATH/bin` is in your PATH:
```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

### Build Manually

```bash
cd /path/to/lazycommands
make build
```

The binary will be created at `bin/lazycommands`.

## Usage

LazyCommands supports two input methods: command-line arguments or stdin.

### Method 1: Command-line Arguments

Run multiple commands in sequence:

```bash
lazycommands 'echo "Starting"' 'sleep 2' 'echo "Done"'
```

### Method 2: Stdin (Heredoc/EOF)

Use heredoc syntax for multi-line command input:

```bash
lazycommands << EOF
echo "Starting"
sleep 2
echo "Done"
EOF
```

Or pipe from a file:

```bash
cat commands.txt | lazycommands
```

Or pipe from echo:

```bash
echo -e "npm install\nnpm run build\nnpm test" | lazycommands
```

### Real-world Examples

**Rails Development Workflow**
```bash
lazycommands 'cd PackageTracker' 'bundle install' 'rails db:migrate' 'rails db:seed'
```

**Node.js Project Setup**
```bash
lazycommands 'npm install' 'npm run build' 'npm test'
```

**Docker Compose Setup**
```bash
lazycommands 'docker-compose down' 'docker-compose pull' 'docker-compose up -d'
```

**Multiple Directory Operations**
```bash
lazycommands 'cd project1 && git pull' 'cd project2 && git pull' 'cd project3 && git pull'
```

**Using Heredoc for Complex Workflows**
```bash
lazycommands << EOF
echo "Starting deployment..."
git pull origin main
npm install
npm run build
npm run test
docker build -t myapp:latest .
docker push myapp:latest
echo "Deployment complete!"
EOF
```

**Using a Commands File**

Create a file `deploy-commands.txt`:
```
echo "Starting deployment..."
npm run lint
npm run test
npm run build
docker-compose up -d
echo "Deployment complete!"
```

Then run:
```bash
cat deploy-commands.txt | lazycommands
```

## How It Works

1. **Sequential Execution**: Commands run one after another in the order specified
2. **Progress Display**: The UI shows a list of all commands with status icons:
   - `⏳` Pending (not started yet)
   - `▶️` Running (currently executing)
   - `✅` Completed (finished successfully)
   - `❌` Failed (exited with error)
   - `⊘` Skipped (not executed due to earlier failure)
3. **Persistent Output**: All output remains in your terminal after the program exits - you can scroll back to see what happened
4. **Fail-Fast Behavior**: If a command fails, execution stops immediately
5. **Error Display**: When a command fails, the full output is displayed so you can diagnose the issue
6. **Summary**: After completion, shows a summary of what happened (e.g., "✅ All commands completed successfully (3/3)")

## Keyboard Controls

| Key | Action |
|-----|--------|
| `q` / `Ctrl+C` | Quit anytime |

## Error Handling

When a command fails (exits with non-zero status), LazyCommands will:
1. Stop execution immediately
2. Mark the failed command with ❌
3. Mark all remaining commands as skipped (⊘)
4. Display the failed command's output
5. Exit with code 1

This fail-fast behavior ensures you catch errors early and can fix them before proceeding.

## Development

### Project Structure

```
lazycommands/
├── main.go                         # Entry point
├── internal/
│   ├── app/                        # Bubble Tea application
│   │   ├── model.go               # State model
│   │   ├── update.go              # Update logic
│   │   ├── view.go                # UI rendering
│   │   └── commands.go            # Command helpers
│   ├── executor/                   # Command execution
│   │   ├── executor.go            # Execution engine
│   │   └── command.go             # Command state
│   ├── ui/                         # UI components
│   │   ├── layout.go              # Layout calculations
│   │   ├── styles.go              # Styling
│   │   └── components.go          # Reusable components
│   └── keys/                       # Keyboard handling
│       └── keymap.go              # Key bindings
├── Makefile                        # Build automation
└── README.md                       # This file
```

### Building

```bash
make build           # Build binary to bin/lazycommands
make install         # Install to $GOPATH/bin
make test            # Run tests
make clean           # Remove build artifacts
```

### Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Styling and layout
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components

## Exit Codes

- `0`: All commands executed successfully
- `1`: At least one command failed

## Future Enhancements

- **Parallel execution**: Run multiple commands concurrently with `-parallel` flag
- **Configuration files**: Define command sequences in YAML/JSON
- **Watch mode**: Re-run commands when files change
- **Command dependencies**: Define which commands depend on others

## License

MIT

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
