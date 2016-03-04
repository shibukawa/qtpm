package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"path/filepath"
)

const packageFileName = "qtpackage.toml"
const userPackageFileName = "qtpackage.user.toml"

type PackageConfig struct {
	Name      string   `toml:"name"`
	Author    string   `toml:"author"`
	License   string   `toml:"license"`
	Type      string   `toml:"type"`
	Requires  []string `toml:"requires"`
	QtModules []string `toml:"qtmodules"`
	Version   []int    `toml:"version"`
}

type PackageUserConfig struct {
	QtDir      string   `toml:"qtdir"`
}

func LoadConfig(dir string, traverse bool) (*PackageConfig, string, error) {
	origDir := dir
	dir, err := filepath.Abs(dir)
	if err != nil {
		return nil, "", err
	}
	for {
		filePath := filepath.Join(dir, packageFileName)
		file, err := os.Open(filePath)
		if err == nil {
			config := &PackageConfig{}
			_, err := toml.DecodeReader(file, config)
			if err != nil {
				return nil, "", err
			}
			return config, dir, nil
		}
		parent := filepath.Dir(dir)
		if dir == parent || !traverse {
			break
		}
		dir = parent
	}
	return nil, "", fmt.Errorf("can't find '%s' at %s", packageFileName, origDir)
}

func LoadUserConfig(dir string) (*PackageUserConfig, error) {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	filePath := filepath.Join(dir, userPackageFileName)
	file, err := os.Open(filePath)
	if err == nil {
		config := &PackageUserConfig{}
		_, err := toml.DecodeReader(file, config)
		if err != nil {
			return nil, err
		}
		return config, nil
	}
	return nil, fmt.Errorf("can't find '%s' at %s", userPackageFileName, dir)
}

func (config *PackageConfig) Save(dir string) error {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return err
	}
	file, err := os.Create(filepath.Join(dir, packageFileName))
	if err != nil {
		return err
	}
	encoder := toml.NewEncoder(file)
	return encoder.Encode(config)
}

func UserName() string {
	names := []string{"LOGNAME", "USER", "USERNAME"}
	for _, name := range names {
		value := os.Getenv(name)
		if value != "" {
			return value
		}
	}
	return "(no name)"
}
