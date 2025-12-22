package executor

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Logger interface to avoid circular dependency
type Logger interface {
	LogCommandStart(cmd *Command)
	LogCommandOutput(cmdID int, line string)
	LogCommandEnd(cmd *Command)
}

// CommandStartedMsg is sent when a command starts executing
type CommandStartedMsg struct {
	Index int
}

// CommandCompletedMsg is sent when a command finishes executing
type CommandCompletedMsg struct {
	Index    int
	ExitCode int
	Error    error
	NewDir   string // If non-empty, working directory changed
}

// TickMsg is sent periodically to refresh the UI and show streaming output
type TickMsg time.Time

// ExecuteCommand runs a command and returns a tea.Cmd that streams output
func ExecuteCommand(index int, cmd *Command, workingDir string, logger Logger) tea.Cmd {
	return func() tea.Msg {
		// Store working directory
		cmd.WorkingDir = workingDir

		// Record start time
		cmd.StartTime = time.Now()
		cmd.Status = StatusRunning

		// Log command start
		if logger != nil {
			logger.LogCommandStart(cmd)
		}

		// Check if this is a cd command
		isCd, targetDir, err := ParseCdCommand(cmd.Raw)
		cmd.IsCdCommand = isCd

		if isCd {
			// Handle cd command specially
			if err != nil {
				// cd command validation failed
				cmd.Status = StatusFailed
				cmd.Error = err
				cmd.EndTime = time.Now()
				cmd.ExitCode = 1
				outputLine := fmt.Sprintf("Error: %v", err)
				cmd.AppendOutput(outputLine)

				// Log output and end
				if logger != nil {
					logger.LogCommandOutput(cmd.ID, outputLine)
					logger.LogCommandEnd(cmd)
				}

				return CommandCompletedMsg{
					Index:    index,
					ExitCode: 1,
					Error:    err,
					NewDir:   "", // No directory change on error
				}
			}

			// cd command succeeded
			cmd.Status = StatusCompleted
			cmd.EndTime = time.Now()
			cmd.ExitCode = 0
			outputLine := fmt.Sprintf("Changed directory to: %s", targetDir)
			cmd.AppendOutput(outputLine)

			// Log output and end
			if logger != nil {
				logger.LogCommandOutput(cmd.ID, outputLine)
				logger.LogCommandEnd(cmd)
			}

			return CommandCompletedMsg{
				Index:    index,
				ExitCode: 0,
				Error:    nil,
				NewDir:   targetDir, // Signal directory change
			}
		}

		// Regular command execution
		// Use the user's shell to support aliases and other shell features
		shell := os.Getenv("SHELL")
		if shell == "" {
			shell = "/bin/sh"
		}

		// Build command that sources shell config to load aliases
		cmdString := cmd.Raw

		// For bash/zsh, prepend source command and use eval to expand aliases
		// Also append pwd output to capture directory changes (including from cd aliases)
		if strings.Contains(shell, "bash") {
			// Source .bashrc if it exists and use eval to expand aliases
			cmdString = "[ -f ~/.bashrc ] && source ~/.bashrc; eval " + shellQuote(cmdString) + "; echo \"__LAZYCOMMANDS_PWD__:$PWD\""
		} else if strings.Contains(shell, "zsh") {
			// Source .zshrc if it exists and use eval to expand aliases
			cmdString = "[ -f ~/.zshrc ] && source ~/.zshrc; eval " + shellQuote(cmdString) + "; echo \"__LAZYCOMMANDS_PWD__:$PWD\""
		}

		execCmd := exec.CommandContext(cmd.ctx, shell, "-c", cmdString)

		// Set working directory if specified
		if workingDir != "" {
			execCmd.Dir = workingDir
		}

		// Get pipes for stdout and stderr
		stdoutPipe, err := execCmd.StdoutPipe()
		if err != nil {
			cmd.Status = StatusFailed
			cmd.Error = err
			cmd.EndTime = time.Now()
			cmd.ExitCode = -1

			// Log command end
			if logger != nil {
				logger.LogCommandEnd(cmd)
			}

			return CommandCompletedMsg{
				Index:    index,
				ExitCode: -1,
				Error:    err,
				NewDir:   "",
			}
		}

		stderrPipe, err := execCmd.StderrPipe()
		if err != nil {
			cmd.Status = StatusFailed
			cmd.Error = err
			cmd.EndTime = time.Now()
			cmd.ExitCode = -1

			// Log command end
			if logger != nil {
				logger.LogCommandEnd(cmd)
			}

			return CommandCompletedMsg{
				Index:    index,
				ExitCode: -1,
				Error:    err,
				NewDir:   "",
			}
		}

		// Start the command
		if err := execCmd.Start(); err != nil {
			cmd.Status = StatusFailed
			cmd.Error = err
			cmd.EndTime = time.Now()
			cmd.ExitCode = -1

			// Log command end
			if logger != nil {
				logger.LogCommandEnd(cmd)
			}

			return CommandCompletedMsg{
				Index:    index,
				ExitCode: -1,
				Error:    err,
				NewDir:   "",
			}
		}

		// Stream output from both stdout and stderr
		var wg sync.WaitGroup
		wg.Add(2)

		go streamOutput(stdoutPipe, cmd, &wg, logger)
		go streamOutput(stderrPipe, cmd, &wg, logger)

		// Wait for output streaming to complete
		wg.Wait()

		// Wait for the command to finish
		err = execCmd.Wait()

		// Record end time
		cmd.EndTime = time.Now()

		// Get exit code
		exitCode := 0
		if err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				exitCode = exitError.ExitCode()
			} else {
				exitCode = -1
			}
			cmd.Status = StatusFailed
			cmd.Error = err
		} else {
			cmd.Status = StatusCompleted
		}

		cmd.ExitCode = exitCode

		// Extract working directory from output if present
		newDir := extractWorkingDir(cmd)

		// Log command end
		if logger != nil {
			logger.LogCommandEnd(cmd)
		}

		return CommandCompletedMsg{
			Index:    index,
			ExitCode: exitCode,
			Error:    err,
			NewDir:   newDir,
		}
	}
}

// streamOutput reads lines from a pipe and appends them to the command's output
func streamOutput(pipe io.ReadCloser, cmd *Command, wg *sync.WaitGroup, logger Logger) {
	defer wg.Done()
	defer pipe.Close()

	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		line := scanner.Text()
		cmd.AppendOutput(line)

		// Log the output line
		if logger != nil {
			logger.LogCommandOutput(cmd.ID, line)
		}
	}
}

// Ticker returns a command that sends TickMsg periodically
func Ticker() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

// shellQuote quotes a string for safe use in shell eval
func shellQuote(s string) string {
	// Replace single quotes with '\'' (end quote, escaped quote, start quote)
	s = strings.ReplaceAll(s, "'", "'\\''")
	return "'" + s + "'"
}

// extractWorkingDir extracts and removes the working directory marker from command output
func extractWorkingDir(cmd *Command) string {
	const marker = "__LAZYCOMMANDS_PWD__:"

	// Check if output has the marker in the last few lines
	if len(cmd.Output) == 0 {
		return ""
	}

	// Look through the last few lines for the marker
	for i := len(cmd.Output) - 1; i >= 0 && i >= len(cmd.Output)-5; i-- {
		line := cmd.Output[i]
		if strings.HasPrefix(line, marker) {
			// Extract the directory path
			dir := strings.TrimPrefix(line, marker)
			dir = strings.TrimSpace(dir)

			// Remove this line from output
			cmd.Output = append(cmd.Output[:i], cmd.Output[i+1:]...)

			return dir
		}
	}

	return ""
}
