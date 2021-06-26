package src

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
)

var keywords map[string]TokenCode

type TTokenRecord struct {
	Token        TokenCode
	TokenString  string
	TokenInteger int32
	TokenFloat   float64
	LineNumber   int
	ColumnNumber int
}

type getTokenFn func() TokenCode

const (
	CR           = rune(13)
	LF           = rune(10)
	EOF_CHAR     = rune(255)
	MAX_INT      = 2147483647
	MAX_EXPONENT = 308
)

type TokenCode byte

const (
	T_EOF TokenCode = iota
	T_STRING
	T_IDENT
	T_INTEGER
	T_FLOAT
	T_PLUS
	T_MINUS
	T_MULT
	T_DIVIDE
	T_LESS
	T_GREATER
	T_LESS_EQ
	T_GREATER_EQ
	T_EQUAL
	T_ASSIGN
	T_NOT_EQ
	T_COLON
	T_SEMICOLON
	T_COMMA
	T_POWER
	T_LPAREN
	T_RPAREN
	T_LBRACKET
	T_RBRACKET
	T_LBRACE
	T_RBRACE
	// keywords
	T_BREAK
	T_IF
	T_DOWNTO
	T_ELSE
	T_THEN
	T_END
	T_TRUE
	T_FALSE
	T_WHILE
	T_DO
	T_REPEAT
	T_UNTIL
	T_FOR
	T_TO
	T_AND
	T_OR
	T_NOT
	T_XOR
	T_DIV
	T_MOD
	T_FUNCTION
	T_REF
	T_RETURN
	T_PRINT
	T_PRINTLN
)

type Scanner struct {
	Token        getTokenFn
	columnNumber int
	lineNumber   int
	ch           rune

	TokenRecord  TTokenRecord
	tokenQueue   []TTokenRecord
	StreamReader *StreamReader

	inMultiLineComment bool
}

func NewScanner() *Scanner {
	s := &Scanner{
		TokenRecord: TTokenRecord{},
		tokenQueue:  []TTokenRecord{},
	}
	s.Token = s.getTokenCode
	s.addKeywords()
	return s
}

func (s *Scanner) addKeywords() {
	keywords = make(map[string]TokenCode)
	keywords["break"] = T_BREAK
	keywords["if"] = T_IF
	keywords["downto"] = T_DOWNTO
	keywords["else"] = T_ELSE
	keywords["then"] = T_THEN
	keywords["end"] = T_END
	keywords["True"] = T_TRUE
	keywords["False"] = T_FALSE
	keywords["while"] = T_WHILE
	keywords["do"] = T_DO
	keywords["repeat"] = T_REPEAT
	keywords["until"] = T_UNTIL
	keywords["for"] = T_FOR
	keywords["to"] = T_TO
	keywords["and"] = T_AND
	keywords["or"] = T_OR
	keywords["not"] = T_NOT
	keywords["xor"] = T_XOR
	keywords["div"] = T_DIV
	keywords["mod"] = T_MOD
	keywords["function"] = T_FUNCTION
	keywords["ref"] = T_REF
	keywords["return"] = T_RETURN
	keywords["print"] = T_PRINT
	keywords["println"] = T_PRINTLN
}

func (s *Scanner) getTokenCode() TokenCode {
	if len(s.tokenQueue) > 0 {
		return s.tokenQueue[0].Token
	}
	return s.TokenRecord.Token
}

func (s *Scanner) ScanString(str string) {
	s.StreamReader = NewStreamReader(str)
	s.startScanner()
}

func (s *Scanner) ScanFile(fileName string) {
	// check for file exist
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		fmt.Printf("the file does not exist: %s", fileName)
	}
	// load file in array of bytes
	fileStream, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("fatal error: could not open the file: %s", fileName)
	}
	s.StreamReader = NewStreamReader(string(fileStream))

	s.startScanner()
}

func (s *Scanner) startScanner() {
	s.lineNumber = 1
	s.columnNumber = 0
	s.ch = s.nextChar()
}

func (s *Scanner) readRawChar() rune {
	if s.StreamReader.EndOfStream() {
		return EOF_CHAR
	}
	s.columnNumber += 1
	return s.StreamReader.Read()
}

func (s *Scanner) getOSIndependentChar() rune {
	ch := s.readRawChar()
	if ch == CR || ch == LF {
		if ch == CR {
			ch = s.readRawChar() // get the LF
			if ch == LF {
				return ch
			} else {
				fmt.Println("expecting line feed character")
				os.Exit(1)
			}
		}
	}
	return ch
}

// retorna el siguiente caracter en el stream de entrada.
// filtra el LineFeed e incrementa el número de línea.
func (s *Scanner) nextChar() rune {
	result := s.getOSIndependentChar()
	// ignorar LF y retornar el siguiente caracter
	if result == LF {
		s.lineNumber += 1
		s.columnNumber = 0
		return rune(' ')
	}
	return result
}

func (s *Scanner) PushBackToken(token TTokenRecord) {
	s.TokenRecord = token
	s.tokenQueue = append(s.tokenQueue, token)
}

func (s *Scanner) NextToken() {
	if len(s.tokenQueue) > 0 {
		s.TokenRecord = s.tokenQueue[0]
		if len(s.tokenQueue) > 0 {
			s.tokenQueue = s.tokenQueue[1:] // discount first element
		}
		return
	}
	s.skipBlanksAndComments()

	// registrar la posición del token que estamos a punto de identificar
	s.TokenRecord.LineNumber = s.lineNumber
	s.TokenRecord.ColumnNumber = s.columnNumber

	if isLetter(s.ch) {
		s.getWord()
		return
	}
	if isDigit(s.ch) || s.ch == rune('.') {
		s.getNumber()
		return
	}
	if s.ch == rune('"') || s.ch == rune('\'') {
		s.getString(s.ch)
		return
	}
	if s.ch == EOF_CHAR {
		s.TokenRecord.Token = T_EOF
		if s.inMultiLineComment {
			fmt.Println("detected unterminated comment, expecting \"*/\"")
			os.Exit(1)
		}
		return
	}
	s.getSpecial()
}

func (s *Scanner) skipBlanksAndComments() {
	for s.ch == rune(' ') || s.ch == rune('\t') || s.ch == rune('/') {
		if s.ch == rune(' ') || s.ch == rune('\t') {
			s.ch = s.nextChar()
		} else {
			// Chequea el inicio del comentario
			if s.StreamReader.Peek() == rune('/') || s.StreamReader.Peek() == rune('*') {
				s.ch = s.getOSIndependentChar()
				if s.ch == rune('/') { // este comentario: // abc - una linea
					s.skipSingleLineComment()
				} else if s.ch == rune('*') {
					s.skipMultiLineComment()
				}
			} else {
				break
			}
		}
	}
}

func (s *Scanner) skipSingleLineComment() {
	for s.ch != LF && s.ch != EOF_CHAR {
		s.ch = s.getOSIndependentChar()
	}
	if s.ch != EOF_CHAR {
		s.ch = s.nextChar() // skip LF
	}
	s.lineNumber += 1
}

// trata con este tipo de comentario: /* ..... */
func (s *Scanner) skipMultiLineComment() {

	s.inMultiLineComment = true
	s.ch = s.nextChar() // skip '*'

	for {
		for s.ch != rune('*') && s.ch != EOF_CHAR {
			s.ch = s.nextChar()
		}
		if s.ch == EOF_CHAR {
			break
		}
		s.ch = s.nextChar() // skip the '*'
		if s.ch == rune('/') {
			s.ch = s.nextChar() // skip the '/'
			s.inMultiLineComment = false
			break
		}
	}
}

func isLetter(ch rune) bool {
	return rune('a') <= ch && ch <= rune('z') || rune('A') <= ch && ch <= rune('Z') || ch == rune('_')
}

func isDigit(ch rune) bool {
	return rune('0') <= ch && ch <= rune('9')
}

func (s *Scanner) getWord() {
	s.TokenRecord.TokenString = ""

	for isLetter(s.ch) || isDigit(s.ch) {
		s.TokenRecord.TokenString += string(s.ch)
		s.ch = s.nextChar()
	}
	s.TokenRecord.Token = isKeyword(s.TokenRecord.TokenString)
}

func isKeyword(tokenString string) TokenCode {
	if keyword, ok := keywords[tokenString]; ok {
		return keyword
	}
	return T_IDENT
}

func (s *Scanner) getNumber() {
	var singleDigit int32
	hasLeftHandSide := false
	hasRightHandSide := false

	s.TokenRecord.TokenInteger = 0
	s.TokenRecord.TokenFloat = 0.0

	// asumimos primero que es un entero
	s.TokenRecord.Token = T_INTEGER
	// chequeamos el punto decimal por si el usuario ha tipeado algo como: .5
	if s.ch != rune('.') {
		hasLeftHandSide = true
		for isDigit(s.ch) {
			singleDigit = int32(s.ch) - int32('0')
			if s.TokenRecord.TokenInteger <= int32((MAX_INT-singleDigit)/10) {
				s.TokenRecord.TokenInteger = 10*s.TokenRecord.TokenInteger + singleDigit
				s.ch = s.nextChar()
			} else {
				fmt.Println("integer Overflow, constant value too large to read")
				os.Exit(1)
			}
		}
	}
	scale := float64(1)
	if s.ch == rune('.') {
		// es un float. Comenzamos coleccionando la parte decimal
		s.TokenRecord.Token = T_FLOAT
		s.TokenRecord.TokenFloat = float64(s.TokenRecord.TokenInteger)
		s.ch = s.nextChar() // skip the period '.'
		if isDigit(s.ch) {
			hasRightHandSide = true
		}

		for isDigit(s.ch) {
			scale *= float64(0.1)
			singleDigit = int32(s.ch) - int32('0')
			s.TokenRecord.TokenFloat = s.TokenRecord.TokenFloat + (float64(singleDigit) * float64(scale))
			s.ch = s.nextChar()
		}
	}
	// revisamos si tenemos un número
	if !hasLeftHandSide && !hasRightHandSide {
		fmt.Println("single period on itw own is not a valid number")
		os.Exit(1)
	}

	exponentSign := 1
	// Chequear la notación cientifica
	if s.ch == rune('e') || s.ch == rune('E') {
		// es un float, comenzamos a coleccionar la parte exponencial
		if s.TokenRecord.Token == T_INTEGER {
			s.TokenRecord.Token = T_FLOAT
			s.TokenRecord.TokenFloat = float64(s.TokenRecord.TokenInteger)
		}
		s.ch = s.nextChar()
		if s.ch == rune('-') || s.ch == rune('+') {
			if s.ch == rune('-') {
				exponentSign = -1
			}
			s.ch = s.nextChar()
		}
		// acumulamos el exponente, revisamos que s.ch sea un digito
		if !isDigit(s.ch) {
			fmt.Println("syntax error: number expected in exponent")
			os.Exit(1)
		}
		evalue := int32(0)
		for isDigit(s.ch) {
			singleDigit = int32(s.ch) - int32('0')
			if evalue <= int32((MAX_EXPONENT-singleDigit)/10) {
				evalue = 10*evalue + singleDigit
				s.ch = s.nextChar()
			} else {
				fmt.Printf("exponent overflow, maximum value for exponent is %d", MAX_EXPONENT)
				os.Exit(1)
			}
		}

		evalue *= int32(exponentSign)
		if s.TokenRecord.Token == T_INTEGER {
			s.TokenRecord.TokenFloat = float64(s.TokenRecord.TokenInteger * int32(math.Pow(10, float64(evalue))))
		} else {
			s.TokenRecord.TokenFloat = s.TokenRecord.TokenFloat * float64(math.Pow(10.0, float64(evalue)))
		}
	}
}

func (s *Scanner) getString(strEnd rune) {
	s.TokenRecord.TokenString = ""
	s.TokenRecord.Token = T_STRING

	s.ch = s.nextChar() // skip the first '"'
	for s.ch != EOF_CHAR {
		if s.ch == rune('\\') {
			s.ch = s.nextChar() // skip the '\'
			switch s.ch {
			case rune('\\'):
				s.TokenRecord.TokenString += string('\\')
			case rune('n'):
				s.TokenRecord.TokenString += string('\n')
			case rune('r'):
				s.TokenRecord.TokenString += string('\r')
			case rune('t'):
				s.TokenRecord.TokenString += string('\t')
			default:
				s.TokenRecord.TokenString += string('\\') + string(s.ch)
			}
			s.ch = s.nextChar()
		} else {
			if s.ch == strEnd {
				s.ch = s.nextChar() // skip the closing '"'
				return
			} else {
				s.TokenRecord.TokenString += string(s.ch)
				s.ch = s.nextChar()
			}
		}
	}
	fmt.Println("string without terminating quotation mark")
	os.Exit(1)
}

func (s *Scanner) getSpecial() {
	switch s.ch {
	case rune('+'):
		s.TokenRecord.Token = T_PLUS
	case rune('-'):
		s.TokenRecord.Token = T_MINUS
	case rune('*'):
		s.TokenRecord.Token = T_MULT
	case rune('/'):
		s.TokenRecord.Token = T_DIVIDE
	case rune('^'):
		s.TokenRecord.Token = T_POWER
	case rune('('):
		s.TokenRecord.Token = T_LPAREN
	case rune(')'):
		s.TokenRecord.Token = T_RPAREN
	case rune('['):
		s.TokenRecord.Token = T_LBRACKET
	case rune(']'):
		s.TokenRecord.Token = T_RBRACKET
	case rune('{'):
		s.TokenRecord.Token = T_LBRACE
	case rune('}'):
		s.TokenRecord.Token = T_RBRACE
	case rune(','):
		s.TokenRecord.Token = T_COMMA
	case rune(';'):
		s.TokenRecord.Token = T_SEMICOLON
	case rune(':'):
		s.TokenRecord.Token = T_COLON
	case rune('<'):
		if s.StreamReader.Peek() == rune('=') {
			s.ch = s.nextChar()
			s.TokenRecord.Token = T_LESS_EQ
		} else {
			s.TokenRecord.Token = T_LESS
		}
	case rune('>'):
		if s.StreamReader.Peek() == rune('=') {
			s.ch = s.nextChar()
			s.TokenRecord.Token = T_GREATER_EQ
		} else {
			s.TokenRecord.Token = T_GREATER
		}
	case rune('!'):
		if s.StreamReader.Peek() == rune('=') {
			s.ch = s.nextChar()
			s.TokenRecord.Token = T_NOT_EQ
		} else {
			fmt.Printf("unexpecting '=' character after explanation point:%v\n", s.ch)
			os.Exit(1)
		}
	case rune('='):
		if s.StreamReader.Peek() == rune('=') {
			s.ch = s.nextChar()
			s.TokenRecord.Token = T_EQUAL
		} else {
			s.TokenRecord.Token = T_ASSIGN
		}
	default:
		fmt.Printf("unrecognized character in source code: %c\n", s.ch)
		os.Exit(1)
	}
	s.ch = s.nextChar()
}

// debug function
func (s *Scanner) TokenToString(tokenCode TokenCode) string {
	switch tokenCode {
	case T_IDENT:
		return fmt.Sprintf("identifier <%s>", s.TokenRecord.TokenString)
	case T_INTEGER:
		return fmt.Sprintf("integer <%d>", s.TokenRecord.TokenInteger)
	case T_FLOAT:
		return fmt.Sprintf("float <%f>", s.TokenRecord.TokenFloat)
	case T_STRING:
		return fmt.Sprintf("string <\"%s\">", s.TokenRecord.TokenString)
	case T_PLUS:
		return fmt.Sprintf("special <'%s'>", "+")
	case T_MINUS:
		return fmt.Sprintf("special <'%s'>", "-")
	case T_MULT:
		return fmt.Sprintf("special <'%s'>", "*")
	case T_DIVIDE:
		return fmt.Sprintf("special <'%s'>", "/")
	case T_POWER:
		return fmt.Sprintf("special <'%s'>", "^")
	case T_LPAREN:
		return fmt.Sprintf("special <'%s'>", "(")
	case T_RPAREN:
		return fmt.Sprintf("special <'%s'>", ")")
	case T_LBRACKET:
		return fmt.Sprintf("special <'%s'>", "[")
	case T_RBRACKET:
		return fmt.Sprintf("special <'%s'>", "]")
	case T_LBRACE:
		return fmt.Sprintf("special <'%s'>", "(")
	case T_RBRACE:
		return fmt.Sprintf("special <'%s'>", ")")
	case T_ASSIGN:
		return fmt.Sprintf("special <'%s'>", "=")
	case T_EQUAL:
		return fmt.Sprintf("special <'%s'>", "==")
	case T_NOT_EQ:
		return fmt.Sprintf("special <'%s'>", "<>")
	case T_LESS:
		return fmt.Sprintf("special <'%s'>", "<")
	case T_LESS_EQ:
		return fmt.Sprintf("special <'%s'>", "<=")
	case T_GREATER:
		return fmt.Sprintf("special <'%s'>", ">")
	case T_GREATER_EQ:
		return fmt.Sprintf("special <'%s'>", ">=")
	case T_COMMA:
		return fmt.Sprintf("special <'%s'>", ",")
	case T_SEMICOLON:
		return fmt.Sprintf("special <'%s'>", ";")
	case T_COLON:
		return fmt.Sprintf("special <'%s'>", ":")
	case T_BREAK:
		return fmt.Sprintf("key word: <'%s'>", s.TokenRecord.TokenString)
	case T_IF:
		return fmt.Sprintf("key word: <'%s'>", s.TokenRecord.TokenString)
	case T_DOWNTO:
		return fmt.Sprintf("key word: <'%s'>", s.TokenRecord.TokenString)
	case T_ELSE:
		return fmt.Sprintf("key word: <'%s'>", s.TokenRecord.TokenString)
	case T_THEN:
		return fmt.Sprintf("key word: <'%s'>", s.TokenRecord.TokenString)
	case T_END:
		return fmt.Sprintf("key word: <'%s'>", s.TokenRecord.TokenString)
	case T_TRUE:
		return fmt.Sprintf("key word: <'%s'>", s.TokenRecord.TokenString)
	case T_FALSE:
		return fmt.Sprintf("key word: <'%s'>", s.TokenRecord.TokenString)
	case T_WHILE:
		return fmt.Sprintf("key word: <'%s'>", s.TokenRecord.TokenString)
	case T_DO:
		return fmt.Sprintf("key word: <'%s'>", s.TokenRecord.TokenString)
	case T_REPEAT:
		return fmt.Sprintf("key word: <'%s'>", s.TokenRecord.TokenString)
	case T_UNTIL:
		return fmt.Sprintf("key word: <'%s'>", s.TokenRecord.TokenString)
	case T_FOR:
		return fmt.Sprintf("key word: <'%s'>", s.TokenRecord.TokenString)
	case T_TO:
		return fmt.Sprintf("key word: <'%s'>", s.TokenRecord.TokenString)
	case T_AND:
		return fmt.Sprintf("key word: <'%s'>", s.TokenRecord.TokenString)
	case T_OR:
		return fmt.Sprintf("key word: <'%s'>", s.TokenRecord.TokenString)
	case T_NOT:
		return fmt.Sprintf("key word: <'%s'>", s.TokenRecord.TokenString)
	case T_XOR:
		return fmt.Sprintf("key word: <'%s'>", s.TokenRecord.TokenString)
	case T_DIV:
		return fmt.Sprintf("key word: <'%s'>", s.TokenRecord.TokenString)
	case T_FUNCTION:
		return fmt.Sprintf("key word: <'%s'>", s.TokenRecord.TokenString)
	case T_REF:
		return fmt.Sprintf("key word: <'%s'>", s.TokenRecord.TokenString)
	case T_PRINT:
		return fmt.Sprintf("key word: <'%s'>", s.TokenRecord.TokenString)
	case T_PRINTLN:
		return fmt.Sprintf("key word: <'%s'>", s.TokenRecord.TokenString)
	}
	return fmt.Sprint("end of stream: <EOF>")
}
