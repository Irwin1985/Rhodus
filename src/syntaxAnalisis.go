package src

import (
	"fmt"
	"os"
)

type SyntaxAnalisis struct {
	sc *Scanner
}

func NewSyntaxAnalisis(sc *Scanner) *SyntaxAnalisis {
	sy := &SyntaxAnalisis{
		sc: sc,
	}
	return sy
}

// expression ::= term { ('+' | '-') term }
func (sy *SyntaxAnalisis) expression() {
	sy.term()
	for sy.sc.Token() == T_PLUS || sy.sc.Token() == T_MINUS {
		sy.sc.NextToken()
		sy.term()
	}
}

// factor ::= '(' expression ')' | number | variable
func (sy *SyntaxAnalisis) factor() {
	switch sy.sc.Token() {
	case T_INTEGER:
		sy.sc.NextToken()
	case T_FLOAT:
		sy.sc.NextToken()
	case T_IDENT:
		sy.sc.NextToken()
	case T_LPAREN:
		sy.sc.NextToken()
		sy.expression()
		sy.expect(T_RPAREN)
	default:
		fmt.Println("expecting identifier, scalar or left parentheses")
	}
}

// term ::= factor { ('+' | '-') factor }
func (sy *SyntaxAnalisis) term() {
	sy.factor()
	for sy.sc.Token() == T_MULT || sy.sc.Token() == T_DIVIDE {
		sy.sc.NextToken()
		sy.factor()
	}
}

// assignment ::= variable '=' expression
func (sy *SyntaxAnalisis) assignment() {
	sy.expect(T_IDENT)
	sy.expect(T_ASSIGN)
	sy.expression()
}

// outputStatement ::= 'println' '(' expression ')'
func (sy *SyntaxAnalisis) outputStatement() {
	sy.sc.NextToken()
	sy.expression()
}

// program ::= assignment | outputStatement
func (sy *SyntaxAnalisis) Program() {
	sy.sc.NextToken()
	switch sy.sc.Token() {
	case T_IDENT:
		sy.assignment()
	case T_PRINT:
		sy.outputStatement()
	default:
		fmt.Println("expecting assignment or print statement")
		os.Exit(1)
	}
}

func (sy *SyntaxAnalisis) expect(tokenCode TokenCode) {
	if tokenCode == sy.sc.getTokenCode() {
		sy.sc.NextToken()
	} else {
		fmt.Printf("expecting:%s\n", sy.sc.TokenToString(tokenCode))
		os.Exit(1)
	}
}
