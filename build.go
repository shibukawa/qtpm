package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func Build() {
	config, dir, err := LoadConfig(".", true)
	if err != nil {
		log.Fatalln(err)
	}
	buildPath := filepath.Join(dir, "build")
	vendorPath := filepath.Join(dir, "vendor")
	if config.Type == "library" {
		AddCMakeForLib(dir, config)
	} else {
		AddCMakeForApp(dir, config)
	}
	os.MkdirAll(buildPath, 0777)
	cmd := exec.Command("cmake", "..")
	cmd.Dir = buildPath
	cmd.Env = append(os.Environ(), "INSTALL_CMAKE_DIR="+vendorPath)
	qtDir := FindQt(dir)
	if qtDir != "" {
		cmd.Env = append(cmd.Env, "QTDIR="+qtDir)
	}
	out, err := cmd.CombinedOutput()
	log.Println(string(out))
	if err != nil {
		log.Fatal(err)
	}

	makeCmd := exec.Command("make")
	makeCmd.Dir = buildPath
	out, err = makeCmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(out))
}

func FindQt(dir string) string {
	env := os.Getenv("QTDIR")
	if env != "" {
		return env
	}
	userSetting, _ := LoadUserConfig(dir)
	if userSetting != nil && userSetting.QtDir != "" {
		return userSetting.QtDir
	}
	var paths []string
	if runtime.GOOS == "windows" {
		paths = []string{
			os.Getenv("USERPROFILE"),
			"C:\\", "D:\\",
			os.Getenv("ProgramFiles"),
			os.Getenv("ProgramFiles(x86)"),
		}
	} else {
		paths = []string{
			os.Getenv("HOME"),
			"/",
		}
	}
	for _, path := range paths {
		versions, err := ioutil.ReadDir(filepath.Join(path, "Qt"))
		if err != nil {
			continue
		}
		var biggestDir string
		for _, version := range versions {
			if strings.HasPrefix(version.Name(), "5.") {
				if version.Name() > biggestDir {
					biggestDir = version.Name()
				}
			}
		}
		if biggestDir == "" {
			continue
		}
		targets, err := ioutil.ReadDir(filepath.Join(path, "Qt", biggestDir))
		for _, target := range targets {
			name := target.Name()
			if strings.Contains(name, "ios") || strings.Contains(name, "android") || strings.Contains(name, "winphone") {
				continue
			}
			return filepath.Join(path, "Qt", biggestDir, name)
		}
	}
	return ""
}
