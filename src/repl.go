package src

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

const (
	RHODUS_VERSION = "1.0"
)

func Repl() {
	var sourceCode string
	var sy *SyntaxAnalisisCalc
	sc := NewScanner()
	sy = NewSyntaxAnalisisCalc(sc)
	scanner := bufio.NewScanner(os.Stdin)
	displayWelcome()
	for {
		displayPrompt()
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		sourceCode = scanner.Text()
		if sourceCode == "quit" {
			break
		}
		sc.ScanString(sourceCode)
		sy.Statement()
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
