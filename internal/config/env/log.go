package env

import (
	"fmt"
	"os"
	"strconv"

	"github.com/pkg/errors"

	"github.com/polshe-v/microservices_auth/internal/config"
)

const (
	logDirectoryEnvName = "LOG_DIRECTORY"
	logFilenameEnvName  = "LOG_FILENAME"
	logMaxSizeEnvName   = "LOG_MAX_SIZE"
	logMaxFilesEnvName  = "LOG_MAX_FILES"
	logMaxAgeEnvName    = "LOG_MAX_AGE"
	logLevelEnvName     = "LOG_LEVEL"

	defaultLogDirectory = "logs/"
	defaultLogFilename  = "auth_service.log"
	defaultLogMaxSize   = 10 // MB
	defaultLogMaxFiles  = 3
	defaultLogMaxAge    = 7 // days
	defaultLogLevel     = "warn"
)

type logConfig struct {
	logDirectory string
	logFilename  string
	logMaxSize   int
	logMaxFiles  int
	logMaxAge    int
	logLevel     string
}

var _ config.LogConfig = (*logConfig)(nil)

// NewLogConfig creates new object of LogConfig interface.
func NewLogConfig() (config.LogConfig, error) {
	logDirectory := os.Getenv(logDirectoryEnvName)
	if len(logDirectory) == 0 {
		logDirectory = defaultLogDirectory
	}

	logFilename := os.Getenv(logFilenameEnvName)
	if len(logFilename) == 0 {
		logFilename = defaultLogFilename
	}

	var logMaxSize int
	logMaxSizeStr := os.Getenv(logMaxSizeEnvName)
	if len(logMaxSizeStr) == 0 {
		logMaxSize = defaultLogMaxSize
	} else {
		res, err := strconv.ParseUint(logMaxSizeStr, 10, 32)
		if err != nil {
			return nil, errors.Errorf("failed to process %s setting", logMaxSizeEnvName)
		}
		logMaxSize = int(res)
	}

	var logMaxFiles int
	logMaxFilesStr := os.Getenv(logMaxFilesEnvName)
	if len(logMaxFilesStr) == 0 {
		logMaxFiles = defaultLogMaxFiles
	} else {
		res, err := strconv.ParseUint(logMaxFilesStr, 10, 32)
		if err != nil {
			return nil, errors.Errorf("failed to process %s setting", logMaxFilesEnvName)
		}
		logMaxFiles = int(res)
	}

	var logMaxAge int
	logMaxAgeStr := os.Getenv(logMaxAgeEnvName)
	if len(logMaxAgeStr) == 0 {
		logMaxAge = defaultLogMaxAge
	} else {
		res, err := strconv.ParseUint(logMaxAgeStr, 10, 32)
		if err != nil {
			return nil, errors.Errorf("failed to process %s setting", logMaxAgeEnvName)
		}
		logMaxAge = int(res)
	}

	logLevel := os.Getenv(logLevelEnvName)
	if len(logLevel) == 0 {
		logLevel = defaultLogLevel
	}

	return &logConfig{
		logDirectory: logDirectory,
		logFilename:  logFilename,
		logMaxSize:   logMaxSize,
		logMaxFiles:  logMaxFiles,
		logMaxAge:    logMaxAge,
		logLevel:     logLevel,
	}, nil
}

func (cfg *logConfig) LogDirectory() string {
	return cfg.logDirectory
}

func (cfg *logConfig) LogFilename() string {
	return cfg.logFilename
}

func (cfg *logConfig) LogFilePath() string {
	return fmt.Sprintf("%s%s", cfg.logDirectory, cfg.logFilename)
}

func (cfg *logConfig) LogMaxSize() int {
	return cfg.logMaxSize
}

func (cfg *logConfig) LogMaxFiles() int {
	return cfg.logMaxFiles
}

func (cfg *logConfig) LogMaxAge() int {
	return cfg.logMaxAge
}

func (cfg *logConfig) LogLevel() string {
	return cfg.logLevel
}
