package main

import (
	"Rhodus/src"
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	fileName := `c:\a1\test.txt`
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		fmt.Println("File not found")
		os.Exit(1)
	}
	fileContent, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	printContent(string(fileContent))

	//repl := src.NewRepl()
	//repl.Start("debug")
	//src.Repl()
	//testCalculator()
	//testScannerWithRepl()
	//testScannerWithFileName()
	//fmt.Println(src.GetSampleScriptsDir())
}

func printContent(content string) {
	r := bufio.NewReader(strings.NewReader(content))
	for {
		if chr, _, err := r.ReadRune(); err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		} else {
			fmt.Printf("%q\n", chr)
		}
	}
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
	sy := src.NewSyntaxAnalisis(sc)
	sy.Program()

	// sc.NextToken()
	// for sc.Token() != src.T_EOF {
	// 	fmt.Println(sc.TokenToString(sc.Token()))
	// 	sc.NextToken()
	// }
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
