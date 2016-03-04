package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"os"
)

var (
	app             = kingpin.New("qtpm", "Package Manager fot Qt")
	verbose         = app.Flag("verbose", "Set verbose mode").Short('v').Bool()
	buildCommand    = app.Command("build", "Build program")
	cleanCommand    = app.Command("clean", "Clean temp files")
	getCommand      = app.Command("get", "Get package")
	saveFlag        = getCommand.Flag("save", "Save package as a dependency module").Bool()
	saveDevFlag     = getCommand.Flag("save-dev", "Save package as a dependency module").Bool()
	getPackageName  = getCommand.Arg("package", "Package name on git repository").String()
	installCommand  = app.Command("install", "Install program")
	testCommand     = app.Command("test", "Test package")
	initCommand     = app.Command("init", "Initialize package")
	initAppCommand  = initCommand.Command("app", "Initialize application")
	appName         = initAppCommand.Arg("name", "Application name").Required().String()
	appLicense      = initAppCommand.Arg("license", "License name").Required().String()
	initLibCommand  = initCommand.Command("lib", "Initialize shared library")
	libName         = initLibCommand.Arg("name", "Library name").Required().String()
	libLicense      = initLibCommand.Arg("license", "License name").Required().String()
	addCommand      = app.Command("add", "Add source template")
	addClassCommand = addCommand.Command("class", "Add class template")
	className       = addClassCommand.Arg("name", "Class name").Required().String()
	addTestCommand  = addCommand.Command("test", "Add test template")
	testName        = addTestCommand.Arg("test", "Test class name").Required().String()
)

func main() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case buildCommand.FullCommand():
		Build()
	case cleanCommand.FullCommand():
		panic("not implemented yet")
	case getCommand.FullCommand():
		panic("not implemented yet: " + *getPackageName)
	case installCommand.FullCommand():
		panic("not implemented yet")
	case testCommand.FullCommand():
		panic("not implemented yet")
	case initAppCommand.FullCommand():
		InitApplication(*appName, *appLicense)
	case initLibCommand.FullCommand():
		InitLibrary(*libName, *libLicense)
	case addClassCommand.FullCommand():
		config, dir, err := LoadConfig(".", true)
		if err != nil {
			log.Fatalln(err)
		}
		AddClass(dir, *className, config.Type == "library")
	case addTestCommand.FullCommand():
		_, dir, err := LoadConfig(".", true)
		if err != nil {
			log.Fatalln(err)
		}
		AddTest(dir, *testName)
	}
}
