package executor

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
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
		// Create the exec command using sh -c to support shell features
		execCmd := exec.CommandContext(cmd.ctx, "sh", "-c", cmd.Raw)

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

		// Log command end
		if logger != nil {
			logger.LogCommandEnd(cmd)
		}

		return CommandCompletedMsg{
			Index:    index,
			ExitCode: exitCode,
			Error:    err,
			NewDir:   "", // No directory change for regular commands
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
