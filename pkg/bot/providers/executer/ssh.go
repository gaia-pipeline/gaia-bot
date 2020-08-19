package executer

import (
	"bytes"
	"context"
	"io/ioutil"
	"strings"

	"github.com/rs/zerolog"
	"golang.org/x/crypto/ssh"

	"github.com/gaia-pipeline/gaia-bot/pkg/bot/providers"
)

// Config has the configuration options for the commenter
type Config struct {
	Address        string
	Port           string
	Username       string
	SSHKeyLocation string
}

// Dependencies defines the dependencies of this commenter
type Dependencies struct {
	Logger    zerolog.Logger
	Converter providers.EnvironmentConverter
}

// SSHExec is an executioner which can execute commands on a remote machine.
type SSHExec struct {
	Config
	Dependencies
}

// NewSSHExecutioner creates a new SSHExec
func NewSSHExecutioner(cfg Config, deps Dependencies) *SSHExec {
	return &SSHExec{
		Config:       cfg,
		Dependencies: deps,
	}
}

// Executioner defines an execution method.
func (e *SSHExec) Execute(ctx context.Context, script string, replace map[string]string) error {
	// Connect to machine via SSH and feed in the command.
	log := e.Logger.With().Interface("replace", replace).Logger()
	log.Debug().Msg("Estabilish SSH session...")

	key, err := ioutil.ReadFile(e.SSHKeyLocation)
	if err != nil {
		log.Error().Err(err).Msg("Unable to read private key")
		return err
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Error().Err(err).Msg("unable to parse private key")
		return err
	}
	username, err := e.Converter.LoadValueFromFile(e.Username)
	if err != nil {
		log.Error().Err(err).Msg("Failed to load address")
		return err
	}
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			// Use the PublicKeys method for remote authentication.
			ssh.PublicKeys(signer),
		},
		// TODO: ssh.FixedHostKey(hostKey),
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	address, err := e.Converter.LoadValueFromFile(e.Address)
	if err != nil {
		log.Error().Err(err).Msg("Failed to load address")
		return err
	}
	client, err := ssh.Dial("tcp", address+":"+e.Port, config)
	if err != nil {
		log.Error().Err(err).Msg("Unable to connect to build server.")
		return err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		log.Error().Err(err).Msg("Failed to create session to server.")
		return err
	}
	defer session.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	// replace placeholders
	for k, v := range replace {
		script = strings.ReplaceAll(script, k, v)
	}

	if err := session.Run(script); err != nil {
		log.Debug().Str("stdout", stdout.String()).Msg("Output of the command...")
		log.Debug().Str("stderr", stderr.String()).Msg("Error of the command...")
		log.Error().Err(err).Msg("Failed to remote run pusher.")
		return err
	}
	return nil
}
