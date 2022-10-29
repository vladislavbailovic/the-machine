package debug

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"the-machine/machine/internal"
	"the-machine/machine/memory"
)

type Dumpable interface {
	Out(byte, *bufio.Writer) (int, error)
	Dump(memory.MemoryAccess) error
	Load() ([]byte, error)
}

type Dumper struct {
	fname string
}

func NewDumper() Dumpable {
	return Dumper{fname: "out.bin"}
}

func (x Dumper) Dump(mem memory.MemoryAccess) error {
	return dumpRawMemory(x, mem, x.fname)
}

func (x Dumper) Load() ([]byte, error) {
	buffer, err := os.ReadFile(x.fname)
	if err != nil {
		return buffer, internal.Error(fmt.Sprintf("error loading dump file %s", x.fname), err, internal.ErrorLoading)
	}
	return buffer, nil
}

func (x Dumper) Out(b byte, w *bufio.Writer) (int, error) {
	return w.Write([]byte{b})
}

type AsciiDumper struct {
	Dumper
	formatter  Formatter
	byteFormat string
}

func NewAsciiDumper(numbers Representation) Dumpable {
	var dump Dumper = Dumper{fname: "out.asc"}
	format := Formatter{
		Numbers:  numbers,
		OutputAs: Byte,
	}
	_, valueFmt := format.GetFormat()
	return AsciiDumper{Dumper: dump, formatter: format, byteFormat: valueFmt}
}

func (x AsciiDumper) Out(b byte, w *bufio.Writer) (int, error) {
	out := " " + fmt.Sprintf(x.byteFormat, b) + " "
	return w.WriteString(out)
}

func (x AsciiDumper) Dump(mem memory.MemoryAccess) error {
	return dumpRawMemory(x, mem, x.fname)
}

func (x AsciiDumper) Load() ([]byte, error) {
	buffer, err := os.ReadFile(x.fname)
	if err != nil {
		return buffer, internal.Error(fmt.Sprintf("error loading dump file %s", x.fname), err, internal.ErrorLoading)
	}

	ascii := strings.Split(strings.TrimSpace(string(buffer)), " ")
	out := []byte{}
	for idx, raw := range ascii {
		if "" == raw {
			continue
		}
		if b, err := strconv.Atoi(raw); err != nil {
			return out, internal.Error(fmt.Sprintf("error loading %s at position %d: %v", x.fname, idx, raw), err, internal.ErrorLoading)
		} else {
			out = append(out, byte(b))
		}
	}
	return out, nil
}

func dumpRawMemory(x Dumpable, mem memory.MemoryAccess, toFname string) error {
	f, err := os.Create(toFname)
	if err != nil {
		return internal.Error(fmt.Sprintf("error creating dump file %s", toFname), err, internal.ErrorSaving)
	}
	defer f.Close()

	buffer := bufio.NewWriter(f)
	idx := 0
	for true {
		if b, err := mem.GetByte(memory.Address(idx)); err != nil {
			break
		} else {
			if _, err := x.Out(b, buffer); err != nil {
				return internal.Error(fmt.Sprintf("error dumping memory to %s", toFname), err, internal.ErrorSaving)
			}
		}
		idx++
	}

	return buffer.Flush()
}
