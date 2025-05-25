package main

import (
	"flag"
	"fmt"

	"github.com/psyame/deej/pkg/deej"
)

var (
	gitCommit  string
	versionTag string
	buildType  string

	verbose bool

	enableLogFile bool
	logPath string

	userConfigPath string
)

func init() {
	flag.BoolVar(&verbose, "verbose", false, "show verbose logs (useful for debugging serial)")
	flag.BoolVar(&verbose, "v", false, "shorthand for --verbose")

	flag.BoolVar(&enableLogFile, "enableLogging", false, "enable output of a log file")
	flag.StringVar(&logPath, "logPath", ".", "the path to the folder in which the log file will be created, will create the directory structure if it doesn't exist. defaults to the current directory when the binary is ran")

	flag.StringVar(&userConfigPath, "config", ".", "the path to the directory containing config.yaml, defaults to looking in the current directory when the binary is ran")
	flag.StringVar(&userConfigPath, "c", ".", "shorthand for --config")
	flag.Parse()
}

func main() {

	// first we need a logger
	logger, err := deej.NewLogger(buildType)
	if err != nil {
		panic(fmt.Sprintf("Failed to create logger: %v", err))
	}

	named := logger.Named("main")
	named.Debug("Created logger")

	named.Infow("Version info",
		"gitCommit", gitCommit,
		"versionTag", versionTag,
		"buildType", buildType)

	// provide a fair warning if the user's running in verbose mode
	if verbose {
		named.Debug("Verbose flag provided, all log messages will be shown")
	}

	// create the deej instance
	d, err := deej.NewDeej(logger, verbose, userConfigPath)
	if err != nil {
		named.Fatalw("Failed to create deej object", "error", err)
	}

	// if injected by build process, set version info to show up in the tray
	if buildType != "" && (versionTag != "" || gitCommit != "") {
		identifier := gitCommit
		if versionTag != "" {
			identifier = versionTag
		}

		versionString := fmt.Sprintf("Version %s-%s", buildType, identifier)
		d.SetVersion(versionString)
	}

	// onwards, to glory
	if err = d.Initialize(); err != nil {
		named.Fatalw("Failed to initialize deej", "error", err)
	}
}
