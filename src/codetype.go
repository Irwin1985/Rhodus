package src

type TByteCode struct {
	OpCode OpCode
	index  int
}

type TProgram []TByteCode
