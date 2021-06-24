package main

import (
	"Rhodus/src"
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	src.Repl()
	//testCalculator()
	//testScannerWithRepl()
	//testScannerWithFileName()
}

func testScannerWithFileName() {
	fileName := `c:\a1\test1.rh`
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		fmt.Printf("file does not exist: %s\n", fileName)
		return
	}
	fmt.Printf("Lexical analysis of file: %s\n", fileName)
	fmt.Println("Test file contents:")
	fmt.Println("-------------------------------")

	fileContents, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("fatal error: could not open the file: %s", fileName)
	}

	fmt.Println(string(fileContents))
	fmt.Println("-------------------------------")

	sc := src.NewScanner()
	sc.ScanString(string(fileContents))

	sc.NextToken()
	for sc.Token() != src.T_EOF {
		fmt.Println(sc.TokenToString(sc.Token()))
		sc.NextToken()
	}
	fmt.Println("\n Success!")
}

func testScannerWithRepl() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Rhodus Scanner Test\nPlease type any valid command and hit enter.\n")
	for {
		fmt.Print(">> ")
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		input := scanner.Text()
		sc := src.NewScanner()
		sc.ScanString(input)
		sc.NextToken()
		for sc.Token() != src.T_EOF {
			fmt.Println(sc.TokenToString(sc.Token()))
			sc.NextToken()
		}
	}
}

func testCalculator() {
	input := `a`
	sc := src.NewScanner()
	sc.ScanString(input)
	sy := src.NewSyntaxAnalisisCalc(sc)
	sy.Statement()
}
