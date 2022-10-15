package instruction

type Type byte

const (
	NOP         Type = 0
	MOV_LIT_AC  Type = iota
	MOV_LIT_R1  Type = iota
	MOV_LIT_R2  Type = iota
	MOV_LIT_R3  Type = iota
	MOV_LIT_R4  Type = iota
	ADD_REG_REG Type = iota
	JNE         Type = iota
	HALT        Type = iota
	END         Type = iota
)
