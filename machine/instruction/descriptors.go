package instruction

import (
	"the-machine/machine/register"
)

var Descriptors = map[Type]Instruction{
	NOP: {
		Description: "No-op",
		Executor:    Passthrough{},
	},

	// Data: registers

	MOV_LIT_R1: {
		Description: "Move literal to register R1",
		Executor:    Lit2Reg{Target: register.R1},
	},
	MOV_LIT_R2: {
		Description: "Move literal to register R2",
		Executor:    Lit2Reg{Target: register.R2},
	},
	MOV_LIT_R3: {
		Description: "Move literal to register R3",
		Executor:    Lit2Reg{Target: register.R3},
	},
	MOV_LIT_R4: {
		Description: "Move literal to register R4",
		Executor:    Lit2Reg{Target: register.R4},
	},
	MOV_LIT_R5: {
		Description: "Move literal to register R5",
		Executor:    Lit2Reg{Target: register.R5},
	},
	MOV_LIT_R6: {
		Description: "Move literal to register R6",
		Executor:    Lit2Reg{Target: register.R6},
	},
	MOV_LIT_R7: {
		Description: "Move literal to register R7",
		Executor:    Lit2Reg{Target: register.R7},
	},
	MOV_LIT_R8: {
		Description: "Move literal to register R8",
		Executor:    Lit2Reg{Target: register.R8},
	},
	MOV_LIT_AC: {
		Description: "Move literal to register Ac",
		Executor:    Lit2Reg{Target: register.Ac},
	},
	MOV_LIT_BNK: {
		Description: "Move literal to register Bnk",
		Executor:    Lit2Reg{Target: register.Bnk},
	},
	MOV_REG_REG: {
		Description: "Copy value from register to register",
		Executor:    Reg2Reg{},
	},

	// Data: memory

	MOV_REG_MEM: {
		Description: "Copy content of register to address in accumulator",
		Executor:    Reg2Mem{},
	},
	MOV_LIT_MEM: {
		Description: "Move literal value to memory address in accumulator",
		Executor:    Lit2Mem{},
	},
	MOV_MEM_REG: {
		Description: "Copy memory at address in register1 to register 2",
		Executor:    Mem2Reg{},
	},

	// Stack

	PUSH_REG: {
		Description: "Push value from register to stack",
		Executor:    Reg2Stack{},
	},
	PUSH_LIT: {
		Description: "Push literal value to stack",
		Executor:    Lit2Stack{},
	},
	POP_REG: {
		Description: "Pop value from stack to register",
		Executor:    Stack2Reg{},
	},

	// Stack math

	ADD_STACK: {
		Description: "Add top 2 stack values and push result",
		Executor:    OperateStack{Operation: OpAdd},
	},
	SUB_STACK: {
		Description: "Subtract second stack value from stack head and push result",
		Executor:    OperateStack{Operation: OpSub},
	},
	MUL_STACK: {
		Description: "Multiply top 2 stack values and push result",
		Executor:    OperateStack{Operation: OpMul},
	},
	DIV_STACK: {
		Description: "Divide stack head by second stack value and push result",
		Executor:    OperateStack{Operation: OpDiv},
	},

	// Math

	ADD_REG_REG: {
		Description: "Add contents of two registers",
		Executor:    OperateReg{Operation: OpAdd},
	},
	ADD_REG_LIT: {
		Description: "Add literal value to register (0-15)",
		Executor:    OperateRegLit{Operation: OpAdd},
	},
	SUB_REG_REG: {
		Description: "Subtract contents of two registers",
		Executor:    OperateReg{Operation: OpSub},
	},
	SUB_REG_LIT: {
		Description: "Sub literal value from register",
		Executor:    OperateRegLit{Operation: OpSub},
	},
	MUL_REG_REG: {
		Description: "Multiply contents of two registers",
		Executor:    OperateReg{Operation: OpMul},
	},
	MUL_REG_LIT: {
		Description: "Multiply register with literal value",
		Executor:    OperateRegLit{Operation: OpMul},
	},
	DIV_REG_REG: {
		Description: "Divide contents of two registers",
		Executor:    OperateReg{Operation: OpDiv},
	},
	DIV_REG_LIT: {
		Description: "Divide register with literal value",
		Executor:    OperateRegLit{Operation: OpDiv},
	},
	MOD_REG_REG: {
		Description: "Remainder of contents of two registers division",
		Executor:    OperateReg{Operation: OpMod},
	},
	MOD_REG_LIT: {
		Description: "Remainder of register with literal value division",
		Executor:    OperateRegLit{Operation: OpMod},
	},

	// Bitwise

	SHL_REG_LIT: {
		Description: "Shift left value in register by literal",
		Executor:    OperateRegLit{Operation: OpShl},
	},
	SHR_REG_LIT: {
		Description: "Shift right value in register by literal",
		Executor:    OperateRegLit{Operation: OpShr},
	},
	AND_REG_LIT: {
		Description: "ANDs value in register by literal",
		Executor:    OperateRegLit{Operation: OpAnd},
	},
	AND_REG_REG: {
		Description: "ANDs value in register by register value",
		Executor:    OperateReg{Operation: OpAnd},
	},
	OR_REG_LIT: {
		Description: "ORs value in register by literal",
		Executor:    OperateRegLit{Operation: OpOr},
	},
	OR_REG_REG: {
		Description: "ORs value in register by register value",
		Executor:    OperateReg{Operation: OpOr},
	},
	XOR_REG_LIT: {
		Description: "XORs value in register by literal",
		Executor:    OperateRegLit{Operation: OpXor},
	},
	XOR_REG_REG: {
		Description: "XORs value in register by register value",
		Executor:    OperateReg{Operation: OpXor},
	},

	// Conditional jumps

	JEQ: {
		Description: "Jump to 2nd register address if Ac equal to 1st parameter register",
		Executor:    Jump{Comparison: CompEq},
	},
	JNE: {
		Description: "Jump to 2nd register address if Ac not equal to 1st parameter register",
		Executor:    Jump{Comparison: CompNe},
	},
	JGT: {
		Description: "Jump to 2nd register address if Ac greater than 1st parameter register",
		Executor:    Jump{Comparison: CompGt},
	},
	JGE: {
		Description: "Jump to 2nd register address if Ac greater than or equal to 1st parameter register",
		Executor:    Jump{Comparison: CompGe},
	},
	JLT: {
		Description: "Jump to 2nd register address if Ac less than 1st parameter register",
		Executor:    Jump{Comparison: CompLt},
	},
	JLE: {
		Description: "Jump to 2nd register address if Ac less than or equal to 1st parameter register",
		Executor:    Jump{Comparison: CompLe},
	},

	// Subroutines

	CALL: {
		Description: "Call subroutine at address in register parameter",
		Executor:    Call{},
	},
	RET: {
		Description: "Return from subroutine",
		Executor:    Return{},
	},
}
