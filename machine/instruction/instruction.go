package instruction

type Instruction byte

const (
	NOP         Instruction = 0
	MOV_LIT_AC  Instruction = iota
	MOV_LIT_R1  Instruction = iota
	MOV_LIT_R2  Instruction = iota
	MOV_LIT_R3  Instruction = iota
	MOV_LIT_R4  Instruction = iota
	ADD_REG_REG Instruction = iota
	JNE         Instruction = iota
	HALT        Instruction = iota
	END         Instruction = iota
)
