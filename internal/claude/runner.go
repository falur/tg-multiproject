package claude

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
)

type RunConfig struct {
	Prompt    string
	CWD       string
	Mode      string // "plan" or "edit"
	Binary    string
	SessionID string // for --resume
}

func Run(ctx context.Context, cfg RunConfig) (<-chan StreamEvent, <-chan error) {
	events := make(chan StreamEvent, 16)
	errc := make(chan error, 1)

	args := buildArgs(cfg)

	cmd := exec.CommandContext(ctx, cfg.Binary, args...)
	cmd.Dir = cfg.CWD

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		errc <- fmt.Errorf("stdout pipe: %w", err)
		close(events)
		close(errc)
		return events, errc
	}

	if err := cmd.Start(); err != nil {
		errc <- fmt.Errorf("start: %w", err)
		close(events)
		close(errc)
		return events, errc
	}

	go func() {
		defer close(events)
		defer close(errc)

		scanner := bufio.NewScanner(stdout)
		scanner.Buffer(make([]byte, 0, 1024*1024), 1024*1024)

		for scanner.Scan() {
			line := scanner.Bytes()
			if len(line) == 0 {
				continue
			}
			var ev StreamEvent
			if err := json.Unmarshal(line, &ev); err != nil {
				continue
			}
			select {
			case events <- ev:
			case <-ctx.Done():
				return
			}
		}

		if err := cmd.Wait(); err != nil {
			if ctx.Err() != nil {
				return
			}
			errc <- fmt.Errorf("process exited: %w", err)
		}
	}()

	return events, errc
}

func buildArgs(cfg RunConfig) []string {
	args := []string{"-p", cfg.Prompt, "--output-format", "stream-json"}

	if cfg.SessionID != "" {
		args = append(args, "--resume", cfg.SessionID)
	}

	switch cfg.Mode {
	case "plan":
		args = append(args, "--permission-mode", "plan")
	case "edit":
		args = append(args, "--allowedTools",
			"Read", "Edit", "Write", "Bash", "Glob", "Grep",
		)
	}

	return args
}
