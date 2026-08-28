# LazyCommands

A terminal UI tool for running shell commands sequentially with visual progress tracking.

https://github.com/user-attachments/assets/44122a95-036c-40c5-8e25-02d3f4d62be0

## Install

### Install script (macOS / Linux)

```bash
curl -fsSL https://raw.githubusercontent.com/alameenkhader/lazycommands/main/install.sh | sh
```

This installs the `lazycommands` binary plus an `lc` alias.

### From source

Requires Go 1.21 or later:

```bash
go install github.com/alameenkhader/lazycommands@latest
```

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

## Future Enhancements

- **Parallel execution**: Run multiple commands concurrently with `-parallel` flag
- **Configuration files**: Define command sequences in YAML/JSON
- **Watch mode**: Re-run commands when files change
- **Command dependencies**: Define which commands depend on others

## License

MIT

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
