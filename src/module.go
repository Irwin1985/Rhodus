package src

type Module struct {
	Name        string
	Code        TProgram
	symbolTable *TSymbolTable
}

func NewModule() *Module {
	m := &Module{
		Code: TProgram{},
	}
	return m
}

func (m *Module) ClearCode() {

}
