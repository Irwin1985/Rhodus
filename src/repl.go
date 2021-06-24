package src

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	RHODUS_VERSION = "1.0"
)

type Repl struct {
	sc *Scanner
	sy *SyntaxAnalisis
}

func NewRepl() *Repl {
	repl := &Repl{
		sc: NewScanner(),
	}
	repl.sy = NewSyntaxAnalisis(repl.sc)

	return repl
}

func (r *Repl) Start(mode string) {
	var sourceCode string
	r.sy = NewSyntaxAnalisis(r.sc)
	scanner := bufio.NewScanner(os.Stdin)
	displayWelcome()
	for { // start the loop
		displayPrompt()
		if mode != "debug" {
			scanned := scanner.Scan()
			if !scanned {
				return
			}
			sourceCode = scanner.Text()
		}
		sourceCode = `run c:\a1\test1.rh`
		if sourceCode == "quit" {
			break
		}
		if sourceCode[0:3] == "run" {
			fileName := strings.TrimSpace(sourceCode[3:])
			if _, err := os.Stat(fileName); os.IsNotExist(err) {
				fmt.Println("File not found:" + fileName)
			} else {
				fileContent, err := ioutil.ReadFile(fileName)
				if err != nil {
					panic(err)
				}
				r.runCode(string(fileContent))
			}
			continue
		} else {
			if runCommand(sourceCode) {
				continue
			}
		}
		if sourceCode != "" {
			r.runCode(sourceCode)
		}
	}
}

func displayPrompt() {
	fmt.Print(">> ")
}

func displayWelcome() {
	currentTime := time.Now()
	fmt.Printf("Welcome to Rhodus, Version %s\n", RHODUS_VERSION)
	fmt.Printf("Data and Time: %s\n", currentTime.Format(time.Stamp))
	fmt.Printf("Type quit to exit\n")
}

func GetSampleScriptsDir() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	parent := filepath.Dir(wd)
	return fmt.Sprintf("%s\\SampleScripts", parent)
}

func runCommand(command string) bool {
	sdir := fmt.Sprintf("%s\\", GetSampleScriptsDir())
	result := false
	if command[0:4] == "list" {
		fileName := command[4:]
		if _, err := os.Stat(fileName); os.IsNotExist(err) {
			fmt.Printf("No such file: %s\n", fileName)
		}
		fmt.Println(ioutil.ReadFile(fileName))
		result = true
	} else if command[0:4] == "edit" {
		fileName := strings.TrimSpace(command[4:])
		exec.Command("notepad.exe", fileName)
		result = true
	} else if command[0:3] == "dir" {
		var files []string
		err := filepath.Walk(sdir, func(path string, info os.FileInfo, err error) error {
			files = append(files, path)
			return nil
		})
		if err != nil {
			panic(err)
		}
		for _, file := range files {
			fmt.Println(file)
		}
		result = true
	}
	return result
}

func (r *Repl) runCode(code string) {
	r.sc.ScanString(code)
	r.sc.NextToken() // start the scanner
	if r.sc.Token() != T_EOF {
		r.sy.Program()
	}
	fmt.Println("Success!")
}
