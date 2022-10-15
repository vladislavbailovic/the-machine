package instruction

type Instruction byte

const (
	NOP         Instruction = 0
	MOV_LIT_R1  Instruction = iota
	MOV_LIT_R2  Instruction = iota
	ADD_REG_REG Instruction = iota
)
