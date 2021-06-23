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
	sc := NewScanner()
	var sy *SyntaxAnalisis
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
		sy = NewSyntaxAnalisis(sc)
		sy.Program()
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
