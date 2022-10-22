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
	instruction.MOV_LIT_AC: {
		Description: "Move literal to register AC",
		Parameters: []instruction.Parameter{
			instruction.ParamLiteral,
		},
		Executor: instruction.Lit2Reg{Target: register.Ac},
	},
	instruction.MOV_LIT_R1: {
		Description: "Move literal to register R1",
		Parameters: []instruction.Parameter{
			instruction.ParamLiteral,
		},
		Executor: instruction.Lit2Reg{Target: register.R1},
	},
	instruction.MOV_LIT_R2: {
		Description: "Move literal to register R2",
		Parameters: []instruction.Parameter{
			instruction.ParamLiteral,
		},
		Executor: instruction.Lit2Reg{Target: register.R2},
	},
	instruction.MOV_LIT_R3: {
		Description: "Move literal to register R3",
		Parameters: []instruction.Parameter{
			instruction.ParamLiteral,
		},
		Executor: instruction.Lit2Reg{Target: register.R3},
	},
	instruction.MOV_LIT_R4: {
		Description: "Move literal to register R4",
		Parameters: []instruction.Parameter{
			instruction.ParamLiteral,
		},
		Executor: instruction.Lit2Reg{Target: register.R4},
	},
	instruction.MOV_REG_REG: {
		Description: "Copy value from register to register",
		Parameters:  []instruction.Parameter{},
		Executor:    instruction.Reg2Reg{},
	},
	instruction.MOV_AC_REG: {
		Description: "Copy value from register to register",
		Parameters:  []instruction.Parameter{},
		Executor:    instruction.Ac2Reg{},
	},
	instruction.ADD_REG_REG: {
		Description: "Add contents of two registers",
		Parameters: []instruction.Parameter{
			instruction.ParamRegister,
			instruction.ParamRegister,
		},
		Executor: instruction.OperateReg{Operation: instruction.OpAdd},
	},
	instruction.ADD_REG_LIT: {
		Description: "Add literal value to register",
		Parameters: []instruction.Parameter{
			instruction.ParamRegister,
			instruction.ParamLiteral,
		},
		Executor: instruction.OperateRegLit{Operation: instruction.OpAdd},
	},
	instruction.SUB_REG_REG: {
		Description: "Subtract contents of two registers",
		Parameters: []instruction.Parameter{
			instruction.ParamRegister,
			instruction.ParamRegister,
		},
		Executor: instruction.OperateReg{Operation: instruction.OpSub},
	},
	instruction.SUB_REG_LIT: {
		Description: "Sub literal value from register",
		Parameters: []instruction.Parameter{
			instruction.ParamRegister,
			instruction.ParamLiteral,
		},
		Executor: instruction.OperateRegLit{Operation: instruction.OpSub},
	},
	instruction.MUL_REG_REG: {
		Description: "Multiply contents of two registers",
		Parameters: []instruction.Parameter{
			instruction.ParamRegister,
			instruction.ParamRegister,
		},
		Executor: instruction.OperateReg{Operation: instruction.OpMul},
	},
	instruction.MUL_REG_LIT: {
		Description: "Multiply register with literal value",
		Parameters: []instruction.Parameter{
			instruction.ParamRegister,
			instruction.ParamLiteral,
		},
		Executor: instruction.OperateRegLit{Operation: instruction.OpMul},
	},
	instruction.DIV_REG_REG: {
		Description: "Divide contents of two registers",
		Parameters: []instruction.Parameter{
			instruction.ParamRegister,
			instruction.ParamRegister,
		},
		Executor: instruction.OperateReg{Operation: instruction.OpDiv},
	},
	instruction.DIV_REG_LIT: {
		Description: "Divide register with literal value",
		Parameters: []instruction.Parameter{
			instruction.ParamRegister,
			instruction.ParamLiteral,
		},
		Executor: instruction.OperateRegLit{Operation: instruction.OpDiv},
	},
	instruction.MOD_REG_REG: {
		Description: "Remainder of contents of two registers division",
		Parameters: []instruction.Parameter{
			instruction.ParamRegister,
			instruction.ParamRegister,
		},
		Executor: instruction.OperateReg{Operation: instruction.OpMod},
	},
	instruction.MOD_REG_LIT: {
		Description: "Remainder of register with literal value division",
		Parameters: []instruction.Parameter{
			instruction.ParamRegister,
			instruction.ParamLiteral,
		},
		Executor: instruction.OperateRegLit{Operation: instruction.OpMod},
	},
	instruction.MOV_REG_MEM: {
		Description: "Copy content of register to address in accumulator",
		Parameters: []instruction.Parameter{
			instruction.ParamRegister,
			instruction.ParamAddress,
		},
		Executor: instruction.Reg2Mem{},
	},
	instruction.MOV_LIT_MEM: {
		Description: "Move literal value to memory address in accumulator",
		Parameters: []instruction.Parameter{
			instruction.ParamLiteral,
			instruction.ParamAddress,
		},
		Executor: instruction.Lit2Mem{},
	},

	instruction.JEQ: {
		Description: "Jump to if Ac equal to literal",
		Parameters: []instruction.Parameter{
			instruction.ParamLiteral,
			instruction.ParamAddress,
		},
		Executor: instruction.Jump{Comparison: instruction.CompEq},
	},
	instruction.JNE: {
		Description: "Jump to if Ac not equal to literal",
		Parameters: []instruction.Parameter{
			instruction.ParamLiteral,
			instruction.ParamAddress,
		},
		Executor: instruction.Jump{Comparison: instruction.CompNe},
	},
	instruction.JGT: {
		Description: "Jump to if Ac greater than literal",
		Parameters: []instruction.Parameter{
			instruction.ParamLiteral,
			instruction.ParamAddress,
		},
		Executor: instruction.Jump{Comparison: instruction.CompGt},
	},
	instruction.JGE: {
		Description: "Jump to if Ac greater than or equal to literal",
		Parameters: []instruction.Parameter{
			instruction.ParamLiteral,
			instruction.ParamAddress,
		},
		Executor: instruction.Jump{Comparison: instruction.CompGe},
	},
	instruction.JLT: {
		Description: "Jump to if Ac less than literal",
		Parameters: []instruction.Parameter{
			instruction.ParamLiteral,
			instruction.ParamAddress,
		},
		Executor: instruction.Jump{Comparison: instruction.CompLt},
	},
	instruction.JLE: {
		Description: "Jump to if Ac less than or equal to literal",
		Parameters: []instruction.Parameter{
			instruction.ParamLiteral,
			instruction.ParamAddress,
		},
		Executor: instruction.Jump{Comparison: instruction.CompLe},
	},
}
