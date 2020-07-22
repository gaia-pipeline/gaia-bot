package executioner

import (
	"bytes"
	"context"
	"os"
	"os/exec"

	"github.com/rs/zerolog"
)

// Config has the configuration options for the commenter
type Config struct {
}

// Dependencies defines the dependencies of this commenter
type Dependencies struct {
	Logger zerolog.Logger
}

// LocalExec is an executioner which can execute commands locally.
type LocalExec struct {
	Config
	Dependencies
}

// NewLocalExecutioner creates a new LocalExec
func NewLocalExecutioner(cfg Config, deps Dependencies) *LocalExec {
	return nil
}

// Executioner defines an execution method.
func (e *LocalExec) Execute(ctx context.Context, script string, args map[string]string) error {
	log := e.Logger
	log.Debug().Msg("Executing script...")
	cmd := exec.Command(script)
	cmd.Env = append(cmd.Env, os.Environ()...)
	for k, v := range args {
		cmd.Env = append(cmd.Env, k+"="+v)
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		log.Debug().Str("stdout", stdout.String()).Msg("Output of the command...")
		log.Debug().Str("stderr", stderr.String()).Msg("Error of the command...")
		log.Error().Err(err).Msg("Failed to run pusher.")
		return err
	}
	return nil
}
