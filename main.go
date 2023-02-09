package main

import (
	"os"

	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/exitcode"
	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/bitrise-steplib/steps-restore-spm-cache/step"
)

func main() {
	exitCode := run()
	os.Exit(int(exitCode))
}

func run() exitcode.ExitCode {
	logger := log.NewLogger()
	envRepo := env.NewRepository()
	inputParser := stepconf.NewInputParser(envRepo)
	cmdFactory := command.NewFactory(envRepo)
	cacheStep := step.New(logger, inputParser, envRepo, cmdFactory)

	if err := cacheStep.Run(); err != nil {
		logger.Errorf(err.Error())
		return exitcode.Failure
	}

	return exitcode.Success
}
