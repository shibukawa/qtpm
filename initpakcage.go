package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func prepareDir(name, license, packageType string) *PackageConfig {
	dir, err := filepath.Abs(".")
	if err != nil {
		log.Fatalln("Can't get directory name")
	}

	licenseKey, licenseName, err := NormalizeLicense(license)
	if err != nil {
		log.Fatalln(err)
	}

	config := &PackageConfig{
		Name:    name,
		Author:  UserName(),
		License: licenseName,
		Type:    packageType,
	}
	WriteLicense(licenseKey)
	initDirs(dir)
	return config
}

func initDirs(workDir string) {
	dirs := []string{"src", "include", "test", "build", "vendor"}
	for _, dir := range dirs {
		os.MkdirAll(filepath.Join(workDir, dir), 0777)
	}
}

func InitLibrary(name, license string) {
	packageName, parentName := ParseName(name)
	config := prepareDir(packageName, license, "library")
	config.Save(".")
	variable := &SourceVariable{
		Target: name,
		Parent: parentName,
	}
	AddClass(".", name, true)
	AddTest(".", name)
	WriteTemplate(".", "include", strings.ToLower(name)+"_global.h", "libglobal.h", variable)
	WriteTemplate(".", "", ".gitignore", "dotgitignore", variable)
	WriteTemplate(".", "", "CMakeExtra.txt", "CMakeExtra.txt", variable)
}

func InitApplication(name, license string) {
	packageName, _ := ParseName(name)
	config := prepareDir(name, license, "application")
	config.QtModules = []string{"Widgets"}
	config.Save(".")

	variable := &SourceVariable{
		Target:    packageName,
		QtModules: []string{"Widgets"},
	}
	WriteTemplate(".", "src", "main.cpp", "main.cpp", variable)
	WriteTemplate(".", "", ".gitignore", "dotgitignore", variable)
	WriteTemplate(".", "", "CMakeExtra.txt", "CMakeExtra.txt", variable)
}
