package step

import (
	"fmt"
	"strings"
	"time"

	"github.com/bitrise-io/go-steputils/v2/cache"
	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/log"
)

const stepId = "restore-spm-cache"

// EnvSwiftPackagesPath is exported by `bitrise-build-cache activate xcode`.
const EnvSwiftPackagesPath = "BITRISE_XCODE_SOURCE_PACKAGES_PATH"

// Cache key templates
// OS + Arch: SPM works on Linux too, and Intel/ARM difference is important on macOS
// checksum: Package.resolved is the dependency lockfile, either in the project root (pure Swift project)
// or at project.xcodeproj/project.xcworkspace/xcshareddata/swiftpm/Package.resolved
var keys = []string{
	`{{ .OS }}-{{ .Arch }}-spm-cache-{{ checksum "**/Package.resolved" }}`,
	`{{ .OS }}-{{ .Arch }}-spm-cache-`,
}

// Build Cache for Xcode relocates DerivedData and the SPM checkouts with it, so those archives hold
// different absolute paths and restore replays an archive to its recorded paths. The marker sits
// before "spm-cache" to keep both namespaces clear of each other's prefix fallback.
var xcelerateKeys = []string{
	`{{ .OS }}-{{ .Arch }}-xcelerate-spm-cache-{{ checksum "**/Package.resolved" }}`,
	`{{ .OS }}-{{ .Arch }}-xcelerate-spm-cache-`,
}

func cacheKeys(envRepo env.Repository) []string {
	if strings.TrimSpace(envRepo.Get(EnvSwiftPackagesPath)) != "" {
		return xcelerateKeys
	}

	return keys
}

type Input struct {
	Verbose bool  `env:"verbose,required"`
	Retries int   `env:"retries,required"`
	Timeout int64 `env:"timeout,required"`
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
	activeKeys := cacheKeys(step.envRepo)
	step.logger.Println()
	step.logger.Printf("Cache keys:")
	step.logger.Printf(strings.Join(activeKeys, "\n"))
	step.logger.Println()

	step.logger.EnableDebugLog(input.Verbose)

	restorer := cache.NewRestorer(step.envRepo, step.logger, step.cmdFactory, nil)
	return restorer.Restore(cache.RestoreCacheInput{
		StepId:         stepId,
		Verbose:        input.Verbose,
		Keys:           activeKeys,
		Timeout:        time.Duration(input.Timeout) * time.Second,
		NumFullRetries: input.Retries,
	})
}
