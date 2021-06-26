package src

import (
	"fmt"
	"os"
)

type VM struct {
	stack     TMachineStack
	stackTop  int
	stackSize int
	module    *Module
}

func NewVM(stackSize int) *VM {
	vm := &VM{
		stackSize: stackSize,
		module:    NewModule(),
	}
	return vm
}

func (vm *VM) createStack(size int) {
	vm.stackSize = size
	vm.stackTop = -1
}

func (vm *VM) RunModule(module *Module) {
	vm.module = module
	vm.run(vm.module.Code)
}

func (vm *VM) run(code TProgram) {
	ip := 0
	for {
		switch code[ip].OpCode {
		case oPushi:
			vm.push(code[ip].index)
		case oAdd:
			vm.addOp()
		case oSub:
			vm.subOp()
		case oMult:
			vm.multOp()
		case oDivide:
			vm.divOp()
		case oUmi:
			vm.unaryMinusOp()
		case oPower:
			vm.pop()
		case oHalt:
			break
		default:
			fmt.Println("unknown opcode encountered in virual machine execution loop")
			os.Exit(1)
		}
		ip += 1
	}
}

func (vm *VM) addOp()        {}
func (vm *VM) subOp()        {}
func (vm *VM) multOp()       {}
func (vm *VM) divOp()        {}
func (vm *VM) unaryMinusOp() {}

func (vm *VM) checkStackOverflow() {
	if vm.stackTop == vm.stackSize {
		fmt.Println("Stack overflow error.")
		os.Exit(1)
	}
}

func (vm *VM) pop() TMachineStackRecord {
	if vm.stackTop >= 0 {
		result := vm.stack[vm.stackTop]
		vm.stackTop -= 1
		return result
	}
	fmt.Println("Stack underflow error")
	os.Exit(1)
	return TMachineStackRecord{}
}

func (vm *VM) push(value interface{}) {
	vm.stackTop += 1
	vm.checkStackOverflow()
	switch value := value.(type) {
	case int:
		vm.stack[vm.stackTop].iValue = value
		vm.stack[vm.stackTop].stackType = stInteger
	case float64:
		vm.stack[vm.stackTop].dValue = value
		vm.stack[vm.stackTop].stackType = stDouble
	case bool:
		vm.stack[vm.stackTop].bValue = value
		vm.stack[vm.stackTop].stackType = stBoolean
	case string:
		vm.stack[vm.stackTop].sValue = value
		vm.stack[vm.stackTop].stackType = stString
	}

}
