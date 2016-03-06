package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
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
	QtModules        []string
	Sources          []string
	Headers          []string
	Resources        []string
	Tests            []string
	ExtraTestSources []string
}

func (sv SourceVariable) Hash() string {
	return "abc"
}

func (sv *SourceVariable) SearchFiles(dir string) {
	sources, err := ioutil.ReadDir(filepath.Join(dir, "src"))
	if err == nil {
		var extraTestSources []string
		for _, source := range sources {
			name := source.Name()
			path := "${PROJECT_SOURCE_DIR}/src/" + name
			if strings.HasSuffix(name, ".cpp") && name != "main.cpp" {
				sv.Sources = append(sv.Sources, path)
				extraTestSources = append(extraTestSources, path)
			} else if strings.HasSuffix(name, ".h") {
				sv.Headers = append(sv.Headers, path)
			}
		}
	}

	resources, err := ioutil.ReadDir(filepath.Join(dir, "resource"))
	if err == nil {
		for _, resource := range resources {
			name := resource.Name()
			if strings.HasSuffix(name, ".qrc") {
				sv.Resources = append(sv.Resources, name)
			}
		}
	}
	tests, err := ioutil.ReadDir(filepath.Join(dir, "test"))
	if err == nil {
		for _, test := range tests {
			name := test.Name()
			if !strings.HasSuffix(name, ".cpp") {
				continue
			}
			if strings.HasSuffix(name, "_test.cpp") {
				sv.Tests = append(sv.Tests, name[:len(name)-4])
			} else {
				sv.ExtraTestSources = append(sv.ExtraTestSources, name)
			}
		}
	}
	sort.Strings(sv.Sources)
	sort.Strings(sv.Headers)
	sort.Strings(sv.Resources)
	sort.Strings(sv.Tests)
	sort.Strings(sv.ExtraTestSources)
}

func AddTest(basePath, name string) {
	name, _ = ParseName(name)
	variable := &SourceVariable{
		Target: name,
	}
	WriteTemplate(basePath, "test", strings.ToLower(name)+"_test.cpp", "testclass.cpp", variable)
}

func AddClass(basePath, name string, isLibrary bool) {
	className, parent := ParseName(name)
	variable := &SourceVariable{
		Target:  className,
		Parent:  parent,
		Library: isLibrary,
	}
	WriteTemplate(basePath, "src", strings.ToLower(className)+".h", "classsource.h", variable)
	WriteTemplate(basePath, "src", strings.ToLower(className)+".cpp", "classsource.cpp", variable)
}

func AddLicense(config *PackageConfig, name string) {
	licenseKey, licenseName, err := NormalizeLicense(name)
	if err != nil {
		log.Fatalln(err)
	}
	config.License = licenseName
	WriteLicense(config.Dir, licenseKey)
	config.Save()
}

func AddCMakeForApp(config *PackageConfig) {
	variable := &SourceVariable{
		Target:    CleanName(config.Name),
		QtModules: InsertCore(config.QtModules),
		Library:   true,
	}
	variable.SearchFiles(config.Dir)
	WriteTemplate(config.Dir, "", "CMakeLists.txt", "CMakeListsApp.txt", variable)
}

func AddCMakeForLib(config *PackageConfig) {
	variable := &SourceVariable{
		Target:    CleanName(config.Name),
		QtModules: InsertCore(config.QtModules),
		Library:   true,
	}
	variable.VersionMajor = config.Version[0]
	variable.VersionMinor = config.Version[1]
	variable.VersionPatch = config.Version[2]
	variable.SearchFiles(config.Dir)

	WriteTemplate(config.Dir, "", "CMakeLists.txt", "CMakeListsLib.txt", variable)
	WriteTemplate(config.Dir, "", "CMakeExtra.txt", "CMakeExtra.txt", variable)
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
	os.MkdirAll(filepath.Dir(filePath), 0744)
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

func InsertCore(modules []string) []string {
	found := false
	for _, module := range modules {
		if module == "Core" {
			found = true
		}
	}
	if !found {
		modules = append(modules, "Core")
	}
	return modules
}
