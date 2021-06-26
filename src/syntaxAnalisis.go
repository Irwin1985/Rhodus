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

// statementList ::= statement { ';' statement }
func (sy *SyntaxAnalisis) statementList() {
	sy.statement()
	for sy.sc.Token() == T_SEMICOLON {
		sy.expect(T_SEMICOLON)
		if sy.sc.Token() == T_UNTIL || sy.sc.Token() == T_END || sy.sc.Token() == T_ELSE || sy.sc.Token() == T_EOF {
			break
		}
		sy.statement()
	}
}

// statement ::= assignment | sorStatement | ifStatement | whileStatement | repeatStatement
// 							| returnStatement | breakStatement | functionDef | printlnStatement | endOfStream
func (sy *SyntaxAnalisis) statement() {
	switch sy.sc.Token() {
	case T_IDENT:
		token1 := sy.sc.TokenRecord
		sy.sc.NextToken()
		token2 := sy.sc.TokenRecord
		if token2.Token == T_ASSIGN {
			sy.assignment()
		} else if token2.Token == T_LPAREN {
			sy.sc.PushBackToken(token1)
			sy.sc.PushBackToken(token2)
			sy.factor()
		}
	case T_IF:
		sy.ifStatement()
	case T_FOR:
		sy.forStatement()
	case T_WHILE:
		sy.whileStatement()
	case T_REPEAT:
		sy.repeatStatement()
	case T_RETURN:
		sy.returnStatement()
	case T_BREAK:
		sy.breakStatement()
	case T_FUNCTION:
		sy.functionDef()
	case T_PRINT, T_PRINTLN:
		sy.printlnStatement()
	default:
		fmt.Println("expecting assignment, if, for, while or repeat statement")
		os.Exit(1)
	}
}

// repeatStatement ::= 'repeat' statementList 'until' expression
func (sy *SyntaxAnalisis) repeatStatement() {
	sy.sc.NextToken() // skip T_REPEAT
	sy.statementList()
	sy.expect(T_UNTIL)
	sy.expression()
}

// whileStatement ::= 'while' expression 'do' statementList 'end'
func (sy *SyntaxAnalisis) whileStatement() {
	sy.sc.NextToken() // skip T_WHILE
	sy.expression()
	sy.expect(T_DO)
	sy.statementList()
	sy.expect(T_END)
}

// forStatement ::= 'for' identifier '=' expression
// ('to' | 'downto' ) expression 'do' statementList 'end'
func (sy *SyntaxAnalisis) forStatement() {
	sy.sc.NextToken() // skip the T_FOR
	sy.expect(T_IDENT)
	sy.expect(T_ASSIGN)
	sy.expression()
	if sy.sc.Token() == T_TO || sy.sc.Token() == T_DOWNTO {
		sy.expect(T_TO)
		sy.expression()
		sy.expect(T_DO)
		sy.statementList()
		sy.expect(T_END)
	} else {
		fmt.Println("expecting 'to' or 'downto' in for loop.")
		os.Exit(1)
	}
}

// breakStatement ::= 'break'
func (sy *SyntaxAnalisis) breakStatement() {
	sy.sc.NextToken()
}

// ifStatement ::= 'if' expression 'then' statementList
func (sy *SyntaxAnalisis) ifStatement() {
	sy.sc.NextToken() // skip T_IF
	sy.expression()
	sy.expect(T_THEN)
	sy.statementList()
	sy.ifEnd()
}

// ifEnd ::= 'end' | 'else' statementList 'end'
func (sy *SyntaxAnalisis) ifEnd() {
	if sy.sc.Token() == T_ELSE {
		sy.sc.NextToken() // skip T_ELSE
		sy.statementList()
		sy.expect(T_END)
	} else {
		sy.expect(T_END)
	}
}

// functionDef ::= 'function' identifier [ '(' argumentList ')' ]
func (sy *SyntaxAnalisis) functionDef() {
	sy.sc.NextToken() // skip T_FUNCTION
	sy.expect(T_IDENT)
	if sy.sc.Token() == T_LPAREN {
		sy.sc.NextToken() // skip T_LPAREN
		if sy.sc.Token() != T_RPAREN {
			sy.argumentList()
		}
		sy.expect(T_RPAREN)
	}
	sy.statementList()
	sy.expect(T_END)
}

// argumentList ::= argument { ',' argument }
func (sy *SyntaxAnalisis) argumentList() {
	sy.argument()
	for sy.sc.Token() == T_COMMA {
		sy.sc.NextToken() // skip T_COMMA
		sy.argument()
	}
}

// argument ::= ['ref'] identifier
func (sy *SyntaxAnalisis) argument() {
	if sy.sc.Token() == T_REF {
		sy.sc.NextToken() // skip T_REF
	}
	sy.expect(T_IDENT)
}

// expression ::= simpleExpression | simpreExpression relationalOp simpleExpression
func (sy *SyntaxAnalisis) expression() {
	sy.simpleExpression()
	if sy.sc.Token() == T_LESS || sy.sc.Token() == T_LESS_EQ || sy.sc.Token() == T_GREATER ||
		sy.sc.Token() == T_GREATER_EQ || sy.sc.Token() == T_EQUAL || sy.sc.Token() == T_NOT_EQ {
		sy.sc.NextToken() // skip matched token
		sy.expression()
	}
}

// simpleExpression ::= term { addingOp term }
func (sy *SyntaxAnalisis) simpleExpression() {
	sy.term()
	for sy.addingOp() {
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
		sy.expect(T_IDENT)
		if sy.sc.Token() == T_LBRACKET { // index assignment or expression
			sy.sc.NextToken() // skip the T_LBRACKET
			if sy.sc.Token() != T_RBRACKET {
				sy.expressionList()
			}
			sy.expect(T_RBRACKET)
		} else if sy.sc.Token() == T_LPAREN { // function call
			sy.expect(T_LPAREN)
			if sy.sc.Token() != T_RPAREN {
				sy.expressionList()
			}
			sy.expect(T_RPAREN)
		}
	case T_LPAREN:
		sy.sc.NextToken()
		sy.expression()
		sy.expect(T_RPAREN)
	case T_STRING:
		sy.sc.NextToken() // skip T_STRING
	case T_NOT: // not booleanExpression
		sy.sc.NextToken()
		sy.expression()
	case T_FALSE:
		sy.sc.NextToken()
	case T_TRUE:
		sy.sc.NextToken()
	case T_LBRACE: // lists: {"1", 2, True, False, etc}
		sy.sc.NextToken() // skip T_LBRACE
		if sy.sc.Token() != T_RBRACE {
			sy.doList()
		}
		sy.expect(T_RBRACE)
	default:
		fmt.Println("expecting scalar, identifier or left parentheses")
	}
}

// doList ::= expression {',' expression}
func (sy *SyntaxAnalisis) doList() {
	sy.expression()
	for sy.sc.Token() == T_COMMA {
		sy.sc.NextToken()
		sy.expression()
	}
}

// term ::= power { multiplyOp power }
func (sy *SyntaxAnalisis) term() {
	sy.power()
	for sy.multiplyOp() {
		sy.sc.NextToken()
		sy.power()
	}
}

// power ::= {'+'|'-'} factor ['^' power]
func (sy *SyntaxAnalisis) power() {
	sign := float64(1)

	for sy.sc.Token() == T_PLUS || sy.sc.Token() == T_MINUS {
		if sy.sc.Token() == T_MINUS {
			sign = sign * -1
			sy.sc.NextToken() // eat T_MINUS
		} else {
			sy.sc.NextToken()
		}
	}

	sy.factor()
	if sy.sc.Token() == T_POWER {
		sy.sc.NextToken()
		sy.factor()
	}
}

// assignment ::= variable '=' expression
func (sy *SyntaxAnalisis) assignment() {
	sy.expect(T_IDENT)
	if sy.sc.Token() == T_LBRACKET {
		sy.sc.NextToken() // skip '['
		sy.expressionList()
		sy.expect(T_RBRACKET)
	}
	sy.expect(T_ASSIGN)
	sy.expression()
}

// addingOp ::= '+' | '-' | or | xor
func (sy *SyntaxAnalisis) addingOp() bool {
	addOp := sy.sc.Token() == T_PLUS || sy.sc.Token() == T_MINUS
	logicalOp := sy.sc.Token() == T_OR || sy.sc.Token() == T_XOR
	return addOp || logicalOp
}

// multiplyOp ::= '*' | '/' | and | mod | div
func (sy *SyntaxAnalisis) multiplyOp() bool {
	multOp := sy.sc.Token() == T_MULT || sy.sc.Token() == T_DIVIDE
	logicalOp := sy.sc.Token() == T_AND || sy.sc.Token() == T_MOD || sy.sc.Token() == T_DIV

	return multOp || logicalOp
}

// expressionList ::= expression { ',' expression }
func (sy *SyntaxAnalisis) expressionList() {
	sy.expression()
	for sy.sc.Token() == T_COMMA {
		sy.expect(T_COMMA)
		sy.expression()
	}
}

// returnStatement ::= 'return' expression
func (sy *SyntaxAnalisis) returnStatement() {
	sy.sc.NextToken() // skip T_RETURN
	sy.expression()
}

// printlnStatement ::= 'println' '(' expression ')'
func (sy *SyntaxAnalisis) printlnStatement() {
	sy.sc.NextToken() // skip T_PRINT, T_PRINTLN
	sy.expect(T_LPAREN)
	sy.expression()
	sy.expect(T_RPAREN)
}

// program ::= assignment | outputStatement
func (sy *SyntaxAnalisis) Program() {
	sy.statementList()
	if sy.sc.Token() != T_EOF {
		sy.expect(T_SEMICOLON)
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
