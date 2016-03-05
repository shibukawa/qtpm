package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func prepareProject(name, license, packageType string) (*PackageConfig, string) {
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
	WriteLicense(dir, licenseKey)
	return config, dir
}

func initDirs(workDir string, extraDirs ...string) {
	dirs := []string{"src", "include", "test", "build", "vendor"}
	dirs = append(dirs, extraDirs...)
	for _, dir := range dirs {
		os.MkdirAll(filepath.Join(workDir, dir), 0777)
	}
}

func InitLibrary(name, license string) {
	packageName, parentName := ParseName(name)
	config, dir := prepareProject(packageName, license, "library")
	initDirs(dir)
	config.Save(".")
	variable := &SourceVariable{
		Target: packageName,
		Parent: parentName,
	}
	AddClass(".", packageName, true)
	AddTest(".", packageName)
	WriteTemplate(".", "include", strings.ToLower(packageName)+"_global.h", "libglobal.h", variable)
	WriteTemplate(".", "", ".gitignore", "dotgitignore", variable)
	WriteTemplate(".", "", "CMakeExtra.txt", "CMakeExtra.txt", variable)
}

func InitApplication(name, license string) {
	packageName, _ := ParseName(name)
	config, dir := prepareProject(packageName, license, "application")
	initDirs(dir, "resource")
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
