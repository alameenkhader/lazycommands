package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/alameen/lazycommands/internal/app"
	"github.com/alameen/lazycommands/internal/executor"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	var commands []*executor.Command

	// Check if stdin has data (piped input)
	stat, _ := os.Stdin.Stat()
	hasStdin := (stat.Mode() & os.ModeCharDevice) == 0

	if hasStdin {
		// Read commands from stdin (one per line)
		commands = readCommandsFromStdin()
	} else if len(os.Args) >= 2 {
		// Parse commands from arguments
		commands = make([]*executor.Command, 0, len(os.Args)-1)
		for i, arg := range os.Args[1:] {
			commands = append(commands, executor.NewCommand(i, arg))
		}
	} else {
		// No input provided
		printUsage()
		os.Exit(1)
	}

	if len(commands) == 0 {
		fmt.Println("Error: No commands provided")
		os.Exit(1)
	}

	// Create the Bubble Tea model
	model := app.NewModel(commands)

	// Create the program (no alt screen - keep output in terminal)
	p := tea.NewProgram(model)

	// Run the program
	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf("\nError running program: %v\n", err)
		os.Exit(1)
	}

	// Get the final model and print summary
	if m, ok := finalModel.(app.Model); ok {
		// Close the logger before exit
		m.CloseLogger()

		fmt.Println() // Add spacing after UI
		printSummary(m)
		os.Exit(m.ExitCode())
	}

	os.Exit(0)
}

// printSummary prints a final summary of what happened
func printSummary(m app.Model) {
	completed := 0
	failed := 0
	skipped := 0

	for _, cmd := range m.Commands() {
		switch cmd.Status {
		case executor.StatusCompleted:
			completed++
		case executor.StatusFailed:
			failed++
		case executor.StatusSkipped:
			skipped++
		}
	}

	total := len(m.Commands())

	if failed > 0 {
		fmt.Printf("❌ Execution failed: %d/%d completed, %d failed, %d skipped\n", completed, total, failed, skipped)
	} else {
		fmt.Printf("✅ All commands completed successfully (%d/%d)\n", completed, total)
	}

	// Print log file path if available
	if logPath := m.LoggerPath(); logPath != "" {
		fmt.Printf("\n📝 Debug log available at: %s\n", logPath)
	}
}

// readCommandsFromStdin reads commands from stdin, one per line
func readCommandsFromStdin() []*executor.Command {
	commands := make([]*executor.Command, 0)
	scanner := bufio.NewScanner(os.Stdin)

	i := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Skip empty lines
		if line != "" {
			commands = append(commands, executor.NewCommand(i, line))
			i++
		}
	}

	return commands
}

// printUsage prints the usage information
func printUsage() {
	fmt.Println("Usage: lazycommands 'cmd1' 'cmd2' 'cmd3' ...")
	fmt.Println("   or: echo 'cmd1' | lazycommands")
	fmt.Println("   or: lazycommands << EOF")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  # Using arguments:")
	fmt.Println("  lazycommands 'echo \"Hello\"' 'sleep 2' 'echo \"Done\"'")
	fmt.Println("  lazycommands 'npm install' 'npm run build' 'npm test'")
	fmt.Println()
	fmt.Println("  # Using heredoc (EOF):")
	fmt.Println("  lazycommands << EOF")
	fmt.Println("  echo \"Starting\"")
	fmt.Println("  sleep 2")
	fmt.Println("  echo \"Done\"")
	fmt.Println("  EOF")
	fmt.Println()
	fmt.Println("  # Using pipe:")
	fmt.Println("  cat commands.txt | lazycommands")
	fmt.Println()
}
