package main

import (
	"Rhodus/src"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	testScannerWithFileName()
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
