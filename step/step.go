package step

import (
	"fmt"
	"strings"

	"github.com/bitrise-io/go-steputils/v2/cache"
	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/log"
)

const stepId = "restore-spm-cache"

// Cache key templates
// OS + Arch: SPM works on Linux too, and Intel/ARM difference is important on macOS
// checksum: Package.resolved is the dependency lockfile, either in the project root (pure Swift project)
// or at project.xcodeproj/project.xcworkspace/xcshareddata/swiftpm/Package.resolved
var keys = []string{
	`{{ .OS }}-{{ .Arch }}-spm-cache-{{ checksum "**/Package.resolved" }}`,
	`{{ .OS }}-{{ .Arch }}-spm-cache-`,
}

type Input struct {
	Verbose bool `env:"verbose,required"`
	Retries int  `env:"retries,required"`
}

type RestoreCacheStep struct {
	logger      log.Logger
	inputParser stepconf.InputParser
	envRepo     env.Repository
	cmdFactory  command.Factory
}

func New(
	logger log.Logger,
	inputParser stepconf.InputParser,
	envRepo env.Repository,
	cmdFactory command.Factory,
) RestoreCacheStep {
	return RestoreCacheStep{
		logger:      logger,
		inputParser: inputParser,
		envRepo:     envRepo,
		cmdFactory:  cmdFactory,
	}
}

func (step RestoreCacheStep) Run() error {
	var input Input
	if err := step.inputParser.Parse(&input); err != nil {
		return fmt.Errorf("failed to parse inputs: %w", err)
	}
	stepconf.Print(input)
	step.logger.Println()
	step.logger.Printf("Cache keys:")
	step.logger.Printf(strings.Join(keys, "\n"))
	step.logger.Println()

	step.logger.EnableDebugLog(input.Verbose)

	restorer := cache.NewRestorer(step.envRepo, step.logger, step.cmdFactory, nil)
	return restorer.Restore(cache.RestoreCacheInput{
		StepId:         stepId,
		Verbose:        input.Verbose,
		Keys:           keys,
		NumFullRetries: input.Retries,
	})
}
