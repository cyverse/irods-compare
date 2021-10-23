package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"

	"github.com/cyverse/irods-compare/pkg/commons"
	"golang.org/x/term"
	"gopkg.in/yaml.v2"

	log "github.com/sirupsen/logrus"
)

// inputMissingParams gets user inputs for parameters missing, such as username and password
func inputMissingParams(config *commons.Config, stdinClosed bool) error {
	logger := log.WithFields(log.Fields{
		"package":  "main",
		"function": "inputMissingParams",
	})

	if len(config.User) == 0 {
		if stdinClosed {
			err := fmt.Errorf("User is not set")
			logger.Error(err)
			return err
		}

		fmt.Print("Username: ")
		fmt.Scanln(&config.User)
	}

	if len(config.Password) == 0 {
		if stdinClosed {
			err := fmt.Errorf("Password is not set")
			logger.Error(err)
			return err
		}

		fmt.Print("Password: ")
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		fmt.Print("\n")
		if err != nil {
			logger.WithError(err).Error("failed to read password")
			return err
		}

		config.Password = string(bytePassword)
	}

	return nil
}

// processArguments processes command-line parameters
func processArguments() (*commons.Config, error, bool) {
	logger := log.WithFields(log.Fields{
		"package":  "main",
		"function": "processArguments",
	})

	var version bool
	var help bool
	var configFilePath string

	config := commons.NewDefaultConfig()

	// Parse parameters
	flag.BoolVar(&version, "version", false, "Print client version information")
	flag.BoolVar(&version, "v", false, "Print client version information (shorthand form)")
	flag.BoolVar(&help, "h", false, "Print help")
	flag.StringVar(&configFilePath, "config", "", "Set Config YAML File")
	flag.StringVar(&config.Host, "host", "", "Set iRODS host")
	flag.IntVar(&config.Port, "port", 1247, "Set iRODS port")
	flag.StringVar(&config.Zone, "zone", "", "Set iRODS zone")
	flag.StringVar(&config.User, "user", "", "Set iRODS user")
	flag.StringVar(&config.User, "u", "", "Set iRODS user (shorthand form)")
	flag.StringVar(&config.Password, "password", "", "Set iRODS password")
	flag.StringVar(&config.Password, "p", "", "Set iRODS password (shorthand form)")

	flag.Parse()

	if version {
		info, err := commons.GetVersionJSON()
		if err != nil {
			logger.WithError(err).Error("failed to get client version info")
			return nil, err, true
		}

		fmt.Println(info)
		return nil, nil, true
	}

	if help {
		flag.Usage()
		return nil, nil, true
	}

	log.SetOutput(os.Stderr)

	stdinClosed := false
	if len(configFilePath) > 0 {
		// read config
		configFileAbsPath, err := filepath.Abs(configFilePath)
		if err != nil {
			logger.WithError(err).Errorf("failed to access the local yaml file %s", configFilePath)
			return nil, err, true
		}

		fileinfo, err := os.Stat(configFileAbsPath)
		if err != nil {
			logger.WithError(err).Errorf("failed to access the local yaml file %s", configFileAbsPath)
			return nil, err, true
		}

		if fileinfo.IsDir() {
			logger.WithError(err).Errorf("local yaml file %s is not a file", configFileAbsPath)
			return nil, fmt.Errorf("local yaml file %s is not a file", configFileAbsPath), true
		}

		yamlBytes, err := ioutil.ReadFile(configFileAbsPath)
		if err != nil {
			logger.WithError(err).Errorf("failed to read the local yaml file %s", configFileAbsPath)
			return nil, err, true
		}

		err = yaml.Unmarshal(yamlBytes, &config)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal YAML - %v", err), true
		}
	}

	// positional arguments
	localFilePath := ""
	if flag.NArg() == 0 {
		logger.Errorf("Need a local file/dir path to compare")
		return nil, fmt.Errorf("Need a local file/dir path to compare"), true
	}

	irodsFilePath := ""
	if flag.NArg() == 1 {
		logger.Errorf("Need an iRODS file/dir path to compare")
		return nil, fmt.Errorf("Need an iRODS file/dir path to compare"), true
	}

	lastArgIdx := flag.NArg() - 1
	localFilePath = flag.Arg(lastArgIdx - 1)
	irodsFilePath = flag.Arg(lastArgIdx)

	config.SourcePath = localFilePath
	config.DestinationPath = irodsFilePath

	err := inputMissingParams(config, stdinClosed)
	if err != nil {
		logger.WithError(err).Error("failed to input missing parameters")
		return nil, err, true
	}

	return config, nil, false
}
