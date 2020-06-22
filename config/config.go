package config

import (
	"errors"
	"github.com/spf13/viper"
)

type Configuration struct {
	CsvDir        string
	DirNesting    int
	DirsInDir     int
	FilesInDir    int
	StringsInFile int
	Measurements  int
}

var WrongDirNesting = errors.New("config DirNesting must be more than 1")
var WrongDirsInDir = errors.New("config DirsInDir must be more than 1")
var WrongFilesInDir = errors.New("config FilesInDir must be more than 1")
var WrongStringsInFile = errors.New("config StringsInFile must be more than 1")
var WrongMeasurements = errors.New("config Measurements must be more than 1")

func Get() (*Configuration, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	var configuration Configuration

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(&configuration); err != nil {
		return nil, err
	}

	switch {
	case configuration.DirNesting < 1:
		return nil, WrongDirNesting
	case configuration.DirsInDir < 1:
		return nil, WrongDirsInDir
	case configuration.FilesInDir < 1:
		return nil, WrongFilesInDir
	case configuration.StringsInFile < 1:
		return nil, WrongStringsInFile
	case configuration.Measurements < 1:
		return nil, WrongMeasurements
	}

	return &configuration, nil
}
