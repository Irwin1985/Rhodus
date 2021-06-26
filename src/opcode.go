package src

type OpCode byte

const (
	oPushi OpCode = iota // Push integer constant onto stack
	oAdd
	oSub
	oMult
	oDivide
	oUmi   // Unary minus
	oPower // x^y
	oHalt
)
