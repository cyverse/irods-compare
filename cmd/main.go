package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/cyverse/go-irodsclient/irods/connection"
	"github.com/cyverse/go-irodsclient/irods/fs"
	"github.com/cyverse/go-irodsclient/irods/types"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	log "github.com/sirupsen/logrus"
)

// IRODSFileNotFoundError ...
type IRODSFileNotFoundError struct {
	message string
}

func (e *IRODSFileNotFoundError) Error() string {
	return e.message
}

// IsIRODSFileNotFoundError evaluates if the given error is IRODSFileNotFoundError
func IsIRODSFileNotFoundError(err error) bool {
	if _, ok := err.(*IRODSFileNotFoundError); ok {
		return true
	}

	return false
}

func main() {
	log.SetLevel(log.DebugLevel)

	logger := log.WithFields(log.Fields{
		"package":  "main",
		"function": "main",
	})

	// parse argument
	config, err, exit := processArguments()
	if err != nil {
		logger.WithError(err).Error("failed to process arguments")
		if exit {
			os.Exit(1)
		}
	}
	if exit {
		os.Exit(0)
	}

	// run
	err = config.Validate()
	if err != nil {
		logger.WithError(err).Error("invalid argument")
		os.Exit(1)
	}

	account, err := types.CreateIRODSAccount(config.Host, config.Port, config.User, config.Zone, types.AuthSchemeNative, config.Password, "")
	if err != nil {
		logger.WithError(err).Errorf("failed to create an iRODSAccount to iRODS %s:%d", config.Host, config.Port)
		os.Exit(1)
	}

	conn := connection.NewIRODSConnection(account, 300*time.Second, "irods-compare")
	err = conn.Connect()
	if err != nil {
		logger.WithError(err).Errorf("failed to connect to iRODS %s:%d", config.Host, config.Port)
		os.Exit(1)
	}

	logger.Infof("Checking local file %s", config.SourcePath)
	sourceFiles, err := checkLocal(config.SourcePath)
	if err != nil {
		logger.WithError(err).Error("invalid source path")
		conn.Disconnect()
		os.Exit(1)
	}

	// find dest file
	isDestinationCollection := true
	_, err = fs.GetCollection(conn, config.DestinationPath)
	if err != nil {
		if types.IsFileNotFoundError(err) {
			// not collection
			isDestinationCollection = false
		} else {
			logger.Error(err)
			conn.Disconnect()
			os.Exit(1)
		}
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Path", "Hash", "File Size", "Modified Time", "Consistent"})

	if !config.Colorize {
		logger.Info("Colorize disabled")
		text.DisableColors()
	}

	for sourceFileIdx, sourceFile := range sourceFiles {
		destinationFile, localSize, irodsSize, localTime, irodsTime, localHash, irodsHash, err := getComparison(conn, config.SourcePath, config.DestinationPath, sourceFile, isDestinationCollection)
		if err != nil {
			if IsIRODSFileNotFoundError(err) {
				// irods file not found
				t.AppendRows([]table.Row{
					{
						sourceFileIdx + 1,
						sourceFile,
						"SKIP",
						localSize,
						localTime,
						text.Colors{text.FgRed}.Sprintf("FALSE"),
					},
					{
						"-->",
						destinationFile,
						text.Colors{text.FgRed}.Sprintf("FILE NOT FOUND"),
						text.Colors{text.FgRed}.Sprintf("FILE NOT FOUND"),
						text.Colors{text.FgRed}.Sprintf("FILE NOT FOUND"),
						text.Colors{text.FgRed}.Sprintf("FALSE"),
					},
				})
			}

			logger.Error(err)
		} else {
			consistent := text.Colors{text.FgRed}.Sprintf("FALSE")
			if localHash == irodsHash && localSize == irodsSize {
				consistent = text.Colors{text.FgGreen}.Sprintf("TRUE")
			}

			t.AppendRows([]table.Row{
				{
					sourceFileIdx + 1,
					sourceFile,
					localHash,
					localSize,
					localTime,
					consistent,
				},
				{
					"-->",
					destinationFile,
					irodsHash,
					irodsSize,
					irodsTime,
					consistent,
				},
			}, table.RowConfig{})
		}

		t.AppendSeparator()
	}
	t.Render()

	conn.Disconnect()
	os.Exit(0)
}

func checkLocal(sourcePath string) ([]string, error) {
	sourcePaths := []string{}

	fileinfo, err := os.Stat(sourcePath)
	if err != nil {
		return sourcePaths, err
	}

	if fileinfo.IsDir() {
		// dir
		absSourcePath, err := filepath.Abs(sourcePath)
		if err != nil {
			return sourcePaths, err
		}

		dirents, err := os.ReadDir(absSourcePath)
		if err != nil {
			return sourcePaths, err
		}

		for _, dirent := range dirents {
			direntpath := filepath.Join(absSourcePath, dirent.Name())
			_sourcePaths, err := checkLocal(direntpath)
			if err != nil {
				return sourcePaths, err
			}

			sourcePaths = append(sourcePaths, _sourcePaths...)
		}

	} else {
		// file
		absSourcePath, err := filepath.Abs(sourcePath)
		if err != nil {
			return sourcePaths, err
		}

		sourcePaths = append(sourcePaths, absSourcePath)
	}

	return sourcePaths, nil
}

func calcChecksum(sourcePath string, hashAlg hash.Hash) (string, error) {
	//hashAlg := md5.New()
	f, err := os.Open(sourcePath)
	if err != nil {
		return "", err
	}

	defer f.Close()

	_, err = io.Copy(hashAlg, f)
	if err != nil {
		return "", err
	}

	sumBytes := hashAlg.Sum(nil)
	sumString := hex.EncodeToString(sumBytes)

	return sumString, nil
}

func getComparison(conn *connection.IRODSConnection, srcPath string, destPath string, srcFile string, isDestPathCollection bool) (string, int64, int64, time.Time, time.Time, string, string, error) {
	emptyTime := time.Time{}
	localFileinfo, err := os.Stat(srcFile)
	if err != nil {
		return "", 0, 0, emptyTime, emptyTime, "", "", err
	}

	destinationFile := destPath
	if isDestPathCollection {
		absSourcePath, err := filepath.Abs(srcPath)
		if err != nil {
			return "", 0, 0, emptyTime, emptyTime, "", "", err
		}

		if absSourcePath == srcFile {
			// source input was a file
			// find the file in dest dir
			sourceFileName := filepath.Base(srcFile)
			destinationFile = path.Join(destinationFile, sourceFileName)
		} else {
			// source input was a directory
			// calc relpath frmo source input and find the file in dest dir
			relSourcePath, err := filepath.Rel(absSourcePath, srcFile)
			if err != nil {
				return "", 0, 0, emptyTime, emptyTime, "", "", err
			}

			destinationFile = path.Join(destinationFile, filepath.ToSlash(relSourcePath))
		}
	}

	// check irods file
	// get parent collection
	destinationDir := path.Dir(destinationFile)
	destinationCollection, err := fs.GetCollection(conn, destinationDir)
	if err != nil {
		return destinationFile, localFileinfo.Size(), 0, localFileinfo.ModTime(), emptyTime, "", "", &IRODSFileNotFoundError{
			message: fmt.Sprintf("failed to find parent dir - %s", destinationDir),
		}
	}

	// get obj
	destinationFileName := path.Base(destinationFile)
	dataobject, err := fs.GetDataObjectMasterReplica(conn, destinationCollection, destinationFileName)
	if err != nil {
		return destinationFile, localFileinfo.Size(), 0, localFileinfo.ModTime(), emptyTime, "", "", &IRODSFileNotFoundError{
			message: fmt.Sprintf("failed to find irods file - %s/%s", destinationDir, destinationFileName),
		}
	}

	// calc hash - md5
	md5hash := md5.New()
	md5sum, err := calcChecksum(srcFile, md5hash)
	if err != nil {
		return destinationFile, 0, 0, emptyTime, emptyTime, "", "", err
	}

	dataobjectTime := dataobject.Replicas[0].CreateTime
	if dataobject.Replicas[0].ModifyTime.After(dataobjectTime) {
		dataobjectTime = dataobject.Replicas[0].ModifyTime
	}

	return destinationFile, localFileinfo.Size(), dataobject.Size, localFileinfo.ModTime(), dataobjectTime, md5sum, dataobject.Replicas[0].CheckSum, nil
}
