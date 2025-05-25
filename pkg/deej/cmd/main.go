package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/psyame/deej/pkg/deej"
	"github.com/psyame/deej/pkg/deej/util"
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
	defaultLogPath := "./logs"
	defaultConfigPath := "."

	if util.Linux() {
		xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
		if xdgConfigHome != "" {
		  defaultLogPath = filepath.Join(xdgConfigHome, "deej", "logs") 
		  defaultConfigPath = filepath.Join(xdgConfigHome, "deej")
	  }
	}

	flag.BoolVar(&verbose, "verbose", false, "show verbose logs (useful for debugging serial) (default: false)")
	flag.BoolVar(&verbose, "v", false, "shorthand for --verbose (default: false)")

	flag.BoolVar(&enableLogFile, "enableLogFile", false, "enable output of a log file (default: false)")
	flag.StringVar(&logPath, "logPath", defaultLogPath, "req. enableLogFile, the path to the folder in which the log file will be created, will create the directory structure if it doesn't exist")

	flag.StringVar(&userConfigPath, "config", defaultConfigPath, "the path to the directory containing config.yaml")
	flag.StringVar(&userConfigPath, "c", defaultConfigPath, "shorthand for --config")
	flag.Parse()
}

func main() {

	// first we need a logger (optionally)
	logger, err := deej.NewLogger(buildType, enableLogFile, logPath)
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
