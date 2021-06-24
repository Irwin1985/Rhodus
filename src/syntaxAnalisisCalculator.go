package src

import (
	"fmt"
	"math"
	"os"
)

type SyntaxAnalisisCalc struct {
	sc          *Scanner
	symbolTable TSymbolTable
}

func NewSyntaxAnalisisCalc(sc *Scanner) *SyntaxAnalisisCalc {
	sy := &SyntaxAnalisisCalc{
		sc:          sc,
		symbolTable: make(map[string]float64),
	}
	// add built-in variables
	sy.symbolTable["pi"] = 3.14159265358979
	return sy
}

// expression ::= term { ('+' | '-') term }
func (sy *SyntaxAnalisisCalc) expression() float64 {
	result := sy.term()
	for sy.sc.Token() == T_PLUS || sy.sc.Token() == T_MINUS {
		if sy.sc.Token() == T_PLUS {
			sy.sc.NextToken()
			result += sy.term()
		} else {
			sy.sc.NextToken()
			result -= sy.term()
		}
	}
	return result
}

// factor ::= '(' expression ')' | number | variable
func (sy *SyntaxAnalisisCalc) factor() float64 {
	switch sy.sc.Token() {
	case T_INTEGER:
		sy.sc.NextToken()
		return float64(sy.sc.TokenRecord.TokenInteger)
	case T_FLOAT:
		sy.sc.NextToken()
		return sy.sc.TokenRecord.TokenFloat
	case T_LPAREN:
		sy.sc.NextToken()
		result := sy.expression()
		sy.expect(T_RPAREN)
		return result
	case T_IDENT:
		return sy.getIdentifierValue(sy.sc.TokenRecord.TokenString)
	default:
		fmt.Println("expecting identifier, scalar or left parentheses")
	}
	return 0
}

func (sy *SyntaxAnalisisCalc) getIdentifierValue(name string) float64 {
	sy.sc.NextToken() // skip T_ASSIGN
	if symbol, ok := sy.symbolTable[name]; ok {
		return symbol
	}
	fmt.Printf("symbol: '%s' has no assigned value\n", name)
	os.Exit(1)
	return 0
}

// assignment ::= identifier '=' expresson
func (sy *SyntaxAnalisisCalc) assignment(variableName string) {
	value := sy.expression()
	if variableName == "pi" {
		fmt.Printf("pi is a built-in constant, you cannot redefine it.\n")
		return
	}
	sy.symbolTable[variableName] = value
}

// term ::= factor { ('+' | '-') factor }
func (sy *SyntaxAnalisisCalc) term() float64 {
	result := sy.power()
	for sy.sc.Token() == T_MULT || sy.sc.Token() == T_DIVIDE {
		if sy.sc.Token() == T_MULT {
			sy.sc.NextToken() // eat T_MULT
			result *= sy.power()
		} else {
			sy.sc.NextToken() // eat T_DIVIDE
			result /= sy.power()
		}
	}
	return result
}

// power ::= {'+'|'-'} factor ['^' power]
func (sy *SyntaxAnalisisCalc) power() float64 {
	sign := float64(1)

	for sy.sc.Token() == T_PLUS || sy.sc.Token() == T_MINUS {
		if sy.sc.Token() == T_MINUS {
			sign = sign * -1
			sy.sc.NextToken() // eat T_MINUS
		} else {
			sy.sc.NextToken()
		}
	}

	result := sy.factor()
	if sy.sc.Token() == T_POWER {
		sy.sc.NextToken()
		result = math.Pow(result, sy.power())
	}
	return sign * result
}

func (sy *SyntaxAnalisisCalc) expect(tokenCode TokenCode) {
	if tokenCode == sy.sc.Token() {
		sy.sc.NextToken()
	} else {
		fmt.Printf("expecting:%s\n", sy.sc.TokenToString(tokenCode))
		os.Exit(1)
	}
}

// statement = assignment | expression
func (sy *SyntaxAnalisisCalc) Statement() {
	sy.sc.NextToken() // start the scanner
	token1 := sy.sc.TokenRecord
	sy.sc.NextToken()
	token2 := sy.sc.TokenRecord

	if sy.sc.Token() == T_ASSIGN {
		if token1.Token == T_IDENT {
			sy.sc.NextToken()
			sy.assignment(token1.TokenString)
		} else {
			fmt.Printf("left-hand side of the assignment must be a variable")
			os.Exit(1)
		}
	} else {
		sy.sc.PushBackToken(token1)
		sy.sc.PushBackToken(token2)
		sy.sc.NextToken()
		fmt.Println(sy.expression())
	}
}
