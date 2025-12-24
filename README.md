# LazyCommands

A terminal UI tool for running shell commands sequentially with visual progress tracking.

https://github.com/user-attachments/assets/44122a95-036c-40c5-8e25-02d3f4d62be0

## Prerequisites
- LazyCommands is built with Go. You need Go 1.21 or later installed on your system.

## Install
- Clone this repository
- cd /path/to/lazycomands
- `go install` - This will install the `lazycommands` binary to `$GOPATH/bin` (usually `~/go/bin`). Make sure `$GOPATH/bin` is in your PATH

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
```

## Future Enhancements

- **Parallel execution**: Run multiple commands concurrently with `-parallel` flag
- **Configuration files**: Define command sequences in YAML/JSON
- **Watch mode**: Re-run commands when files change
- **Command dependencies**: Define which commands depend on others

## License

MIT

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
