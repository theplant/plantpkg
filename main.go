package main

import (
	"errors"
	"fmt"
	"go/build"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/manifoldco/promptui"
)

func main() {

	var gopath = build.Default.GOPATH
	if len(gopath) == 0 {
		fmt.Println("GOPATH not set, please set and continue")
		return
	}

	fmt.Println("Your GOPATH:", gopath)

	validate := func(input string) error {
		// _, err := strconv.ParseFloat(input, 64)
		// if err != nil {
		// 	return errors.New("Invalid number")
		// }
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Generate go package",
		Validate: validate,
		Default:  "github.com/theplant/mynewpkg",
	}

	generateGoPackagePath, err := prompt.Run()

	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	pkgSegs := strings.Split(generateGoPackagePath, "/")
	packageName := pkgSegs[len(pkgSegs)-1]

	validateServiceName := func(input string) error {
		if len([]rune(input)) == 0 {
			return errors.New("service name required.")
		}
		if !unicode.IsUpper([]rune(input)[0]) {
			return errors.New("first charactor has to be upper case.")
		}
		return nil
	}

	promptService := promptui.Prompt{
		Label:    "Service Name",
		Validate: validateServiceName,
		Default:  strings.Title(packageName),
	}

	serviceName, err := promptService.Run()

	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	dir := filepath.Join(gopath, "src", generateGoPackagePath)
	_, err = os.Stat(dir)
	if err == nil {
		fmt.Printf("%s already exists, remove it first to generate.\n", dir)
		return
	}

	templatePath := filepath.Join(gopath, "src", "github.com/theplant/plantpkg/template")

	cmd := exec.Command("cp", "-r", templatePath, dir)
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
		return
	}

	replaceInFiles(dir, "github.com/theplant/plantpkg/template", generateGoPackagePath)

	replaceInFiles(dir, "template", packageName)
	replaceInFiles(dir, "Template", serviceName)

	fmt.Printf("You choose %q\n", generateGoPackagePath)
}

func replaceInFiles(baseDir string, from, to string) {
	fileList := []string{}
	err := filepath.Walk(baseDir, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			fileList = append(fileList, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	for _, file := range fileList {
		replaceInFile(file, from, to)
	}
}

func replaceInFile(filepath, from, to string) {
	read, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}

	newContents := strings.Replace(string(read), from, to, -1)

	// fmt.Println(newContents)

	err = ioutil.WriteFile(filepath, []byte(newContents), 0)
	if err != nil {
		panic(err)
	}
}
