package executor

import (
	"bufio"
	"io"
	"os/exec"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// CommandStartedMsg is sent when a command starts executing
type CommandStartedMsg struct {
	Index int
}

// CommandCompletedMsg is sent when a command finishes executing
type CommandCompletedMsg struct {
	Index    int
	ExitCode int
	Error    error
}

// TickMsg is sent periodically to refresh the UI and show streaming output
type TickMsg time.Time

// ExecuteCommand runs a command and returns a tea.Cmd that streams output
func ExecuteCommand(index int, cmd *Command) tea.Cmd {
	return func() tea.Msg {
		// Record start time
		cmd.StartTime = time.Now()
		cmd.Status = StatusRunning

		// Create the exec command using sh -c to support shell features
		execCmd := exec.CommandContext(cmd.ctx, "sh", "-c", cmd.Raw)

		// Get pipes for stdout and stderr
		stdoutPipe, err := execCmd.StdoutPipe()
		if err != nil {
			cmd.Status = StatusFailed
			cmd.Error = err
			cmd.EndTime = time.Now()
			return CommandCompletedMsg{
				Index:    index,
				ExitCode: -1,
				Error:    err,
			}
		}

		stderrPipe, err := execCmd.StderrPipe()
		if err != nil {
			cmd.Status = StatusFailed
			cmd.Error = err
			cmd.EndTime = time.Now()
			return CommandCompletedMsg{
				Index:    index,
				ExitCode: -1,
				Error:    err,
			}
		}

		// Start the command
		if err := execCmd.Start(); err != nil {
			cmd.Status = StatusFailed
			cmd.Error = err
			cmd.EndTime = time.Now()
			return CommandCompletedMsg{
				Index:    index,
				ExitCode: -1,
				Error:    err,
			}
		}

		// Stream output from both stdout and stderr
		var wg sync.WaitGroup
		wg.Add(2)

		go streamOutput(stdoutPipe, cmd, &wg)
		go streamOutput(stderrPipe, cmd, &wg)

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

		return CommandCompletedMsg{
			Index:    index,
			ExitCode: exitCode,
			Error:    err,
		}
	}
}

// streamOutput reads lines from a pipe and appends them to the command's output
func streamOutput(pipe io.ReadCloser, cmd *Command, wg *sync.WaitGroup) {
	defer wg.Done()
	defer pipe.Close()

	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		line := scanner.Text()
		cmd.AppendOutput(line)
	}
}

// Ticker returns a command that sends TickMsg periodically
func Ticker() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}
