package cmd

import (
	"os/exec"
	"strings"

	"github.com/rs/zerolog/log"
)

const (
	Whitespace = " "
)

func RunCommand(command string, projectPath string) {
	splitCmd := strings.Split(command, Whitespace)
	cmd := exec.Command(splitCmd[0], splitCmd[1:]...)
	cmd.Dir = projectPath

	err := cmd.Run()
	if err != nil {
		log.Error().Err(err).Msg("Could not run `" + command + "`.")
	}
}
