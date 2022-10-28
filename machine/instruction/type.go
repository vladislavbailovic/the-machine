package instruction

import "fmt"

// Actually 6 bits = 64 instructions max
type Type byte

const (
	NOP Type = 0

	PUSH_REG Type = iota
	PUSH_LIT Type = iota
	POP_REG  Type = iota

	MOV_LIT_R1 Type = iota
	MOV_LIT_R2 Type = iota
	MOV_LIT_R3 Type = iota
	MOV_LIT_R4 Type = iota
	MOV_LIT_R5 Type = iota
	MOV_LIT_R6 Type = iota
	MOV_LIT_R7 Type = iota
	MOV_LIT_R8 Type = iota

	MOV_REG_REG Type = iota
	MOV_REG_MEM Type = iota
	MOV_LIT_MEM Type = iota
	MOV_MEM_REG Type = iota

	ADD_REG_REG Type = iota
	ADD_REG_LIT Type = iota
	ADD_STACK   Type = iota

	SUB_REG_REG Type = iota
	SUB_REG_LIT Type = iota
	SUB_STACK   Type = iota

	MUL_REG_REG Type = iota
	MUL_REG_LIT Type = iota
	MUL_STACK   Type = iota

	DIV_REG_REG Type = iota
	DIV_REG_LIT Type = iota
	DIV_STACK   Type = iota

	MOD_REG_REG Type = iota
	MOD_REG_LIT Type = iota

	SHL_REG_LIT Type = iota
	SHR_REG_LIT Type = iota

	AND_REG_LIT Type = iota
	AND_REG_REG Type = iota

	OR_REG_LIT Type = iota
	OR_REG_REG Type = iota

	XOR_REG_LIT Type = iota
	XOR_REG_REG Type = iota

	JNE Type = iota
	JEQ Type = iota
	JGT Type = iota
	JGE Type = iota
	JLT Type = iota
	JLE Type = iota

	CALL Type = iota
	RET  Type = iota

	HALT Type = iota

	_sizeofType = iota
)

// Safeguard assertion for number of instruction types
var _compileCheck uint8 = 63 - _sizeofType

func (x Type) AsByte() byte {
	return byte(x)
}

func (x Type) Pack(raw ...uint16) []byte {
	var value uint16
	switch len(raw) {
	case 0:
		value = 0
	case 1:
		value = raw[0]
	case 2:
		v1 := raw[0] << 12
		v2 := raw[1] << 8
		value = (v1 | v2) >> 8
	default:
		panic("can't pack more than 2 bytes worth of data atm")
	}
	instr := value | (uint16(x.AsByte()) << 10) // shift 6 instruction bits
	return []byte{
		byte(instr),
		byte(instr >> 8),
	}
}

func Decode(rawInstruction uint16) (Type, uint16) {
	instructionType := byte(
		((rawInstruction >> 10) & 0b0000_0000_0011_1111), // extract 6 instruction bits
	)
	kind := Type(instructionType)
	raw := rawInstruction & 0b0000_0011_1111_1111 // mask off 6 instruction bits for params remainder
	return kind, raw
}

type unpacker struct{}

// Unpacks individual parameter bytes packed by instruction::Pack
func (x unpacker) unpack(raw uint16) []byte {
	b1 := byte(
		(raw & 0b0000_0000_1111_0000) >> 4,
	)
	b2 := byte(
		raw & 0b0000_0000_0000_1111,
	)
	// fmt.Printf("\t- raw: %016b (%d)\n", raw, raw)
	// fmt.Printf("\t-  b1: %016b (%d)\n", b1, b1)
	// fmt.Printf("\t-  b2: %016b (%d)\n", b2, b2)

	return []byte{b1, b2}
}

type Op byte

const (
	OpAdd Op = 0
	OpSub Op = iota
	OpMul Op = iota
	OpDiv Op = iota
	OpMod Op = iota
	OpShl Op = iota
	OpShr Op = iota
	OpAnd Op = iota
	OpOr  Op = iota
	OpXor Op = iota
)

func (x Op) String() string {
	switch x {
	case OpAdd:
		return "+"
	case OpSub:
		return "-"
	case OpMul:
		return "*"
	case OpDiv:
		return "/"
	case OpMod:
		return "%"
	case OpShl:
		return "<<"
	case OpShr:
		return ">>"
	case OpAnd:
		return "&"
	case OpOr:
		return "|"
	case OpXor:
		return "^"
	}
	return fmt.Sprintf("unknown operator: %d", x)
}

type Comparison byte

const (
	CompNe Comparison = 0
	CompEq Comparison = iota
	CompGe Comparison = iota
	CompGt Comparison = iota
	CompLe Comparison = iota
	CompLt Comparison = iota
)

func (x Comparison) String() string {
	switch x {
	case CompNe:
		return "!="
	case CompEq:
		return "=="
	case CompGt:
		return ">"
	case CompGe:
		return ">="
	case CompLt:
		return "<"
	case CompLe:
		return "<="
	}
	return fmt.Sprintf("unknown comparision: %d", x)
}
