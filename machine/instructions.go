package machine

import (
	"the-machine/machine/instruction"
	"the-machine/machine/register"
)

var Instructions = map[instruction.Type]instruction.Instruction{
	instruction.NOP: {
		Description: "No-op",
		Executor:    instruction.Passthrough{},
	},
	instruction.MOV_LIT_R1: {
		Description: "Move literal to register R1",
		Executor:    instruction.Lit2Reg{Target: register.R1},
	},
	instruction.MOV_LIT_R2: {
		Description: "Move literal to register R2",
		Executor:    instruction.Lit2Reg{Target: register.R2},
	},
	instruction.MOV_LIT_R3: {
		Description: "Move literal to register R3",
		Executor:    instruction.Lit2Reg{Target: register.R3},
	},
	instruction.MOV_LIT_R4: {
		Description: "Move literal to register R4",
		Executor:    instruction.Lit2Reg{Target: register.R4},
	},
	instruction.MOV_LIT_R7: {
		Description: "Move literal to register R7",
		Executor:    instruction.Lit2Reg{Target: register.R7},
	},
	instruction.MOV_REG_REG: {
		Description: "Copy value from register to register",
		Executor:    instruction.Reg2Reg{},
	},
	instruction.ADD_REG_REG: {
		Description: "Add contents of two registers",
		Executor:    instruction.OperateReg{Operation: instruction.OpAdd},
	},
	instruction.ADD_REG_LIT: {
		Description: "Add literal value to register",
		Executor:    instruction.OperateRegLit{Operation: instruction.OpAdd},
	},
	instruction.SUB_REG_REG: {
		Description: "Subtract contents of two registers",
		Executor:    instruction.OperateReg{Operation: instruction.OpSub},
	},
	instruction.SUB_REG_LIT: {
		Description: "Sub literal value from register",
		Executor:    instruction.OperateRegLit{Operation: instruction.OpSub},
	},
	instruction.MUL_REG_REG: {
		Description: "Multiply contents of two registers",
		Executor:    instruction.OperateReg{Operation: instruction.OpMul},
	},
	instruction.MUL_REG_LIT: {
		Description: "Multiply register with literal value",
		Executor:    instruction.OperateRegLit{Operation: instruction.OpMul},
	},
	instruction.DIV_REG_REG: {
		Description: "Divide contents of two registers",
		Executor:    instruction.OperateReg{Operation: instruction.OpDiv},
	},
	instruction.DIV_REG_LIT: {
		Description: "Divide register with literal value",
		Executor:    instruction.OperateRegLit{Operation: instruction.OpDiv},
	},
	instruction.MOD_REG_REG: {
		Description: "Remainder of contents of two registers division",
		Executor:    instruction.OperateReg{Operation: instruction.OpMod},
	},
	instruction.MOD_REG_LIT: {
		Description: "Remainder of register with literal value division",
		Executor:    instruction.OperateRegLit{Operation: instruction.OpMod},
	},
	instruction.MOV_REG_MEM: {
		Description: "Copy content of register to address in accumulator",
		Executor:    instruction.Reg2Mem{},
	},
	instruction.MOV_LIT_MEM: {
		Description: "Move literal value to memory address in accumulator",
		Executor:    instruction.Lit2Mem{},
	},
	instruction.JEQ: {
		Description: "Jump to 2nd register address if Ac equal to 1st parameter register",
		Executor:    instruction.Jump{Comparison: instruction.CompEq},
	},
	instruction.JNE: {
		Description: "Jump to 2nd register address if Ac not equal to 1st parameter register",
		Executor:    instruction.Jump{Comparison: instruction.CompNe},
	},
	instruction.JGT: {
		Description: "Jump to 2nd register address if Ac greater than 1st parameter register",
		Executor:    instruction.Jump{Comparison: instruction.CompGt},
	},
	instruction.JGE: {
		Description: "Jump to 2nd register address if Ac greater than or equal to 1st parameter register",
		Executor:    instruction.Jump{Comparison: instruction.CompGe},
	},
	instruction.JLT: {
		Description: "Jump to 2nd register address if Ac less than 1st parameter register",
		Executor:    instruction.Jump{Comparison: instruction.CompLt},
	},
	instruction.JLE: {
		Description: "Jump to 2nd register address if Ac less than or equal to 1st parameter register",
		Executor:    instruction.Jump{Comparison: instruction.CompLe},
	},
}
