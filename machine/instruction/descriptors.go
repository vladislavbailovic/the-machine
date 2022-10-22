package instruction

import (
	"the-machine/machine/register"
)

var Descriptors = map[Type]Instruction{
	NOP: {
		Description: "No-op",
		Executor:    Passthrough{},
	},
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
	MOV_LIT_R7: {
		Description: "Move literal to register R7",
		Executor:    Lit2Reg{Target: register.R7},
	},
	MOV_REG_REG: {
		Description: "Copy value from register to register",
		Executor:    Reg2Reg{},
	},
	ADD_REG_REG: {
		Description: "Add contents of two registers",
		Executor:    OperateReg{Operation: OpAdd},
	},
	ADD_REG_LIT: {
		Description: "Add literal value to register",
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
	MOV_REG_MEM: {
		Description: "Copy content of register to address in accumulator",
		Executor:    Reg2Mem{},
	},
	MOV_LIT_MEM: {
		Description: "Move literal value to memory address in accumulator",
		Executor:    Lit2Mem{},
	},
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
}
