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
			instruction.PARAM16,
		},
		Executor: instruction.Lit2Reg{Target: register.Ac},
	},
	instruction.MOV_LIT_R1: {
		Description: "Move literal to register R1",
		Parameters: []instruction.Parameter{
			instruction.PARAM16,
		},
		Executor: instruction.Lit2Reg{Target: register.R1},
	},
	instruction.MOV_LIT_R2: {
		Description: "Move literal to register R2",
		Parameters: []instruction.Parameter{
			instruction.PARAM16,
		},
		Executor: instruction.Lit2Reg{Target: register.R2},
	},
	instruction.MOV_LIT_R3: {
		Description: "Move literal to register R3",
		Parameters: []instruction.Parameter{
			instruction.PARAM16,
		},
		Executor: instruction.Lit2Reg{Target: register.R3},
	},
	instruction.MOV_LIT_R4: {
		Description: "Move literal to register R4",
		Parameters: []instruction.Parameter{
			instruction.PARAM16,
		},
		Executor: instruction.Lit2Reg{Target: register.R2},
	},
	instruction.ADD_REG_REG: {
		Description: "Add contents of two registers",
		Parameters: []instruction.Parameter{
			instruction.PARAM8,
			instruction.PARAM8,
		},
		Executor: instruction.OperateReg{Operation: instruction.OpAdd},
	},
	instruction.ADD_REG_LIT: {
		Description: "Add literal value to register",
		Parameters: []instruction.Parameter{
			instruction.PARAM8,
			instruction.PARAM16,
		},
		Executor: instruction.OperateRegLit{Operation: instruction.OpAdd},
	},
	instruction.SUB_REG_REG: {
		Description: "Subtract contents of two registers",
		Parameters: []instruction.Parameter{
			instruction.PARAM8,
			instruction.PARAM8,
		},
		Executor: instruction.OperateReg{Operation: instruction.OpSub},
	},
	instruction.SUB_REG_LIT: {
		Description: "Sub literal value from register",
		Parameters: []instruction.Parameter{
			instruction.PARAM8,
			instruction.PARAM16,
		},
		Executor: instruction.OperateRegLit{Operation: instruction.OpSub},
	},
	instruction.MUL_REG_REG: {
		Description: "Multiply contents of two registers",
		Parameters: []instruction.Parameter{
			instruction.PARAM8,
			instruction.PARAM8,
		},
		Executor: instruction.OperateReg{Operation: instruction.OpMul},
	},
	instruction.MUL_REG_LIT: {
		Description: "Multiply register with literal value",
		Parameters: []instruction.Parameter{
			instruction.PARAM8,
			instruction.PARAM16,
		},
		Executor: instruction.OperateRegLit{Operation: instruction.OpMul},
	},
	instruction.DIV_REG_REG: {
		Description: "Divide contents of two registers",
		Parameters: []instruction.Parameter{
			instruction.PARAM8,
			instruction.PARAM8,
		},
		Executor: instruction.OperateReg{Operation: instruction.OpDiv},
	},
	instruction.DIV_REG_LIT: {
		Description: "Divide register with literal value",
		Parameters: []instruction.Parameter{
			instruction.PARAM8,
			instruction.PARAM16,
		},
		Executor: instruction.OperateRegLit{Operation: instruction.OpDiv},
	},
	instruction.JNE: {
		Description: "Jump if not equal",
		Parameters: []instruction.Parameter{
			instruction.PARAM16,
			instruction.PARAM16,
		},
		Executor: instruction.Jump{},
	},
}
