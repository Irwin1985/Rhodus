package src

type TStackType byte

const (
	stInteger TStackType = iota
	stBoolean
	stDouble
	stString
	stList
)

type TMachineStackRecord struct {
	stackType TStackType
	iValue    int
	bValue    bool
	dValue    float64
	sValue    string
	lValue    interface{}
}

type PMachineStackRecord *TMachineStackRecord
type TMachineStack []TMachineStackRecord
