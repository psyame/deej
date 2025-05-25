package deej

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/psyame/deej/pkg/deej/util"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	buildTypeNone    = ""
	buildTypeDev     = "dev"
	buildTypeRelease = "release"

	logFilename  = "deej-latest-run.log"
)
// logDirectory has to be non-constant as it's set within the NewLogger function, also expose if log file is even enabled
var (
	logDirectory string
	logFileEnabled bool
)


// NewLogger provides a logger instance for the whole program
func NewLogger(buildType string, fileEnabled bool, logPath string) (*zap.SugaredLogger, error) {
	var loggerConfig zap.Config

	// release with file enabled: info and above, log to file only (no UI), otherwise log to stderr
	if buildType == buildTypeRelease {
		loggerConfig = zap.NewProductionConfig()
		loggerConfig.Encoding = "console"
		
		logFileEnabled = fileEnabled

		if fileEnabled {
		  if err := util.EnsureDirExists(logPath); err != nil {
			  return nil, fmt.Errorf("ensure log directory exists: %w", err)
		  }

			// change output to file
			loggerConfig.OutputPaths = []string{filepath.Join(logPath, logFilename)}
			logDirectory = logPath
	  } else {
		  // make it colorful
			loggerConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		}
		// development: debug and above, log to stderr only, colorful
	} else {
		loggerConfig = zap.NewDevelopmentConfig()

		// make it colorful
		loggerConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// all build types: make it readable
	loggerConfig.EncoderConfig.EncodeCaller = nil
	loggerConfig.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}

	loggerConfig.EncoderConfig.EncodeName = func(s string, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(fmt.Sprintf("%-27s", s))
	}

	logger, err := loggerConfig.Build()
	if err != nil {
		return nil, fmt.Errorf("create zap logger: %w", err)
	}

	// no reason not to use the sugared logger - it's fast enough for anything we're gonna do
	sugar := logger.Sugar()

	return sugar, nil
}
