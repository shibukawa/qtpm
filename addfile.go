package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

type SourceVariable struct {
	Target           string
	TargetSmall      string
	TargetLarge      string
	Parent           string
	Library          bool
	VersionMinor     int
	VersionMajor     int
	VersionPatch     int
	HasTest          bool
	QtModules        []string
	Sources          []string
	Headers          []string
	Tests            []string
	ExtraTestSources string
}

func AddTest(basePath, name string) {
	name, _ = ParseName(name)
	variable := &SourceVariable{
		Target: name,
	}
	WriteTemplate(basePath, "test", "test_"+strings.ToLower(name)+".cpp", "testclass.cpp", variable)
}

func AddClass(basePath, name string, isLibrary bool) {
	className, parent := ParseName(name)
	variable := &SourceVariable{
		Target:  className,
		Parent:  parent,
		Library: isLibrary,
	}
	WriteTemplate(basePath, "include", strings.ToLower(className)+".h", "classsource.h", variable)
	WriteTemplate(basePath, "src", strings.ToLower(className)+".cpp", "classsource.cpp", variable)
}

func AddLicense(basePath string, config *PackageConfig, name string) {
	licenseKey, licenseName, err := NormalizeLicense(name)
	if err != nil {
		log.Fatalln(err)
	}
	config.License = licenseName
	WriteLicense(basePath, licenseKey)
	config.Save(basePath)
}

func AddCMakeForApp(basePath string, config *PackageConfig) {
	variable := &SourceVariable{
		Target:    CleanName(config.Name),
		QtModules: config.QtModules,
		Library:   true,
	}
	switch len(config.Version) {

	}
	sources, err := ioutil.ReadDir("src")
	if err != nil {
		return
	}
	var extraTestSources []string
	for _, source := range sources {
		if strings.HasSuffix(source.Name(), ".cpp") && source.Name() != "main.cpp" {
			variable.Sources = append(variable.Sources, "${PROJECT_SOURCE_DIR}/src/"+source.Name())
			extraTestSources = append(extraTestSources, "${PROJECT_SOURCE_DIR}/src/"+source.Name())
		}
	}

	headers, err := ioutil.ReadDir("include")
	if err != nil {
		return
	}
	for _, header := range headers {
		variable.Headers = append(variable.Headers, "${PROJECT_SOURCE_DIR}/include/"+header.Name())
	}

	tests, err := ioutil.ReadDir("test")
	if err == nil {
		for _, test := range tests {
			if !strings.HasSuffix(test.Name(), ".cpp") {
				continue
			}
			if strings.HasPrefix(test.Name(), "test_") {
				variable.Tests = append(variable.Tests, test.Name()[:len(test.Name())-4])
			} else {
				extraTestSources = append(extraTestSources, "${PROJECT_SOURCE_DIR}/test/"+test.Name())
			}
		}
	}
	variable.HasTest = len(variable.Tests) > 0
	variable.ExtraTestSources = strings.Join(extraTestSources, " ")
	WriteTemplate(basePath, "", "CMakeLists.txt", "CMakeListsApp.txt", variable)
}

func AddCMakeForLib(basePath string, config *PackageConfig) {
	variable := &SourceVariable{
		Target:    CleanName(config.Name),
		QtModules: config.QtModules,
		Library:   true,
	}
	switch len(config.Version) {
	case 0:
		variable.VersionMajor = 1
	case 1:
		variable.VersionMajor = config.Version[0]
	case 2:
		variable.VersionMajor = config.Version[0]
		variable.VersionMinor = config.Version[1]
	default:
		variable.VersionMajor = config.Version[0]
		variable.VersionMinor = config.Version[1]
		variable.VersionPatch = config.Version[2]
	}
	sources, err := ioutil.ReadDir("src")
	if err != nil {
		return
	}
	var extraTestSources []string
	for _, source := range sources {
		if strings.HasSuffix(source.Name(), ".cpp") {
			variable.Sources = append(variable.Sources, "${PROJECT_SOURCE_DIR}/src/"+source.Name())
			extraTestSources = append(extraTestSources, "${PROJECT_SOURCE_DIR}/src/"+source.Name())
		}
	}
	tests, err := ioutil.ReadDir("test")
	if err == nil {
		for _, test := range tests {
			if !strings.HasSuffix(test.Name(), ".cpp") {
				continue
			}
			if strings.HasPrefix(test.Name(), "test_") {
				variable.Tests = append(variable.Tests, test.Name()[:len(test.Name())-4])
			}
			extraTestSources = append(extraTestSources, "${PROJECT_SOURCE_DIR}/test/"+test.Name())
		}
	}
	variable.HasTest = len(variable.Tests) > 0
	variable.ExtraTestSources = strings.Join(extraTestSources, " ")

	WriteTemplate(basePath, "", "CMakeLists.txt", "CMakeListsLib.txt", variable)
	WriteTemplate(basePath, "", "CMakeExtra.txt", "CMakeExtra.txt", variable)
}

func WriteTemplate(basePath, dir, fileName, templateName string, variable *SourceVariable) error {
	variable.TargetSmall = strings.ToLower(variable.Target)
	variable.TargetLarge = strings.ToUpper(variable.Target)
	if variable.Parent == "" {
		variable.Parent = "QObject"
	}
	var filePath string
	var err error
	if dir == "" {
		filePath, err = filepath.Abs(filepath.Join(basePath, fileName))
	} else {
		filePath, err = filepath.Abs(filepath.Join(basePath, dir, fileName))
	}
	if err != nil {
		return err
	}
	os.MkdirAll(filepath.Dir(filePath), 0777)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	src := MustAsset("templates/" + templateName)
	tmp := template.Must(template.New(templateName).Delims("[[", "]]").Parse(string(src)))
	err = tmp.Execute(file, variable)
	if err != nil {
		log.Fatalln(err)
	}
	return err
}

func ParseName(name string) (string, string) {
	names := strings.Split(name, "@")
	className := names[0]
	if strings.HasPrefix(className, "Test") {
		className = className[4:]
	} else if className == "" {
		path, _ := filepath.Abs(".")
		_, className = filepath.Split(path)
	}
	className = strings.ToUpper(className[:1]) + className[1:]

	var parentName string
	if len(names) == 2 {
		parentName = strings.ToUpper(names[1][:2]) + names[1][2:]
	}

	return CleanName(className), CleanName(parentName)
}

var re1 = regexp.MustCompile("[^a-zA-Z0-9_-]")
var re2 = regexp.MustCompile("[-]")

func CleanName(name string) string {
	return re2.ReplaceAllString(re1.ReplaceAllString(name, ""), "_")
}
