package machine

import (
	"encoding/binary"
	"testing"
	"the-machine/machine/register"
)

func Test_Execute_Program(t *testing.T) {
	cpu := NewCpu()
	b1 := make([]byte, 2)
	binary.LittleEndian.PutUint16(b1, uint16(1312))
	b2 := make([]byte, 2)
	binary.LittleEndian.PutUint16(b2, uint16(161))
	cpu.LoadProgram(0, []byte{
		byte(MOV_LIT_R1), b1[0], b1[1],
		byte(MOV_LIT_R2), b2[0], b2[1],
		byte(ADD_REG_REG), byte(register.R1), byte(register.R2),
	})

	step := 0
	for step < 20 {
		if err := cpu.tick(); err != nil {
			t.Fatalf("error at tick %d: %v", step, err)
		}
		step++
	}
	cpu.debug()
}
