package debug

import (
	"bufio"
	"os"
	"the-machine/machine/memory"
)

type Dumper struct {
	fname string
}

func NewDumper() Dumper {
	return Dumper{fname: "out.bin"}
}

func (x Dumper) Dump(mem memory.MemoryAccess) error {
	f, err := os.Create(x.fname)
	if err != nil {
		return err
	}
	defer f.Close()

	buffer := bufio.NewWriter(f)
	idx := 0
	for true {
		if b, err := mem.GetByte(memory.Address(idx)); err != nil {
			break
		} else {
			if _, err := buffer.Write([]byte{b}); err != nil {
				return err
			}
		}
		idx++
	}

	return buffer.Flush()
}
