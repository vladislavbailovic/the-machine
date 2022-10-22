package instruction

type Type byte

const (
	NOP Type = 0

	MOV_LIT_R1 Type = iota
	MOV_LIT_R2 Type = iota
	MOV_LIT_R3 Type = iota
	MOV_LIT_R4 Type = iota
	MOV_LIT_R7 Type = iota

	MOV_REG_REG Type = iota
	MOV_REG_MEM Type = iota
	MOV_LIT_MEM Type = iota

	ADD_REG_REG Type = iota
	ADD_REG_LIT Type = iota

	SUB_REG_REG Type = iota
	SUB_REG_LIT Type = iota

	MUL_REG_REG Type = iota
	MUL_REG_LIT Type = iota

	DIV_REG_REG Type = iota
	DIV_REG_LIT Type = iota

	MOD_REG_REG Type = iota
	MOD_REG_LIT Type = iota

	JNE Type = iota
	JEQ Type = iota
	JGT Type = iota
	JGE Type = iota
	JLT Type = iota
	JLE Type = iota

	HALT Type = iota
)

func (x Type) AsByte() byte {
	return byte(x)
}

type Op byte

const (
	OpAdd Op = 0
	OpSub Op = iota
	OpMul Op = iota
	OpDiv Op = iota
	OpMod Op = iota
)

type Comparison byte

const (
	CompNe Comparison = 0
	CompEq Comparison = iota
	CompGe Comparison = iota
	CompGt Comparison = iota
	CompLe Comparison = iota
	CompLt Comparison = iota
)
