package device

import (
	"fmt"
	"io"
	"os"
	"the-machine/machine/internal"
	"the-machine/machine/memory"
)

// Max 255 descriptors
type FileDescriptor byte

const (
	Stdin  FileDescriptor = 0
	Stdout FileDescriptor = iota
	Stderr FileDescriptor = iota

	_fileDescriptorsLimt FileDescriptor = iota
)

func (x FileDescriptor) String() string {
	switch x {
	case Stdin:
		return "STDIN"
	case Stdout:
		return "STDOUT"
	case Stderr:
		return "STDERR"
	default:
		return fmt.Sprintf("FD#%d", x)
	}
}

type AccessType uint8

const (
	Read  AccessType = 1 << iota
	Write AccessType = 1 << iota
)

func (x AccessType) String() string {
	switch x {
	case Read:
		return "Read"
	case Write:
		return "Write"
	default:
		return "Unknown"
	}
}

type filelike struct {
	descriptor FileDescriptor
	access     AccessType
	stream     interface{}
}

func (x filelike) String() string {
	stream := "<stream>"
	if x.stream == nil {
		stream = "<NO STREAM>"
	}
	return fmt.Sprintf("%s <%s>: %s", x.descriptor, x.access, stream)
}

func (x filelike) Read() (byte, error) {
	if x.access != Read {
		return 0, internal.Error(fmt.Sprintf("unable to read file descriptor in %v", x), nil, internal.ErrorLoading)
	}
	if reader, ok := x.stream.(io.Reader); ok {
		buf := make([]byte, 1, 1)
		if n, err := reader.Read(buf); err != nil {
			return 0, internal.Error(fmt.Sprintf("read error: %v", x), err, internal.ErrorLoading)
		} else if n > 0 {
			return buf[0], nil
		}
		return 0, internal.Error(fmt.Sprintf("read error: EOF %v", x), nil, internal.ErrorLoading)
	}
	return 0, internal.Error(fmt.Sprintf("not a reader: %v", x), nil, internal.ErrorLoading)
}

func (x filelike) Write(b byte) error {
	if x.access != Write {
		return internal.Error(fmt.Sprintf("unable to write to file descriptor in %v", x), nil, internal.ErrorLoading)
	}
	if writer, ok := x.stream.(io.Writer); ok {
		if _, err := writer.Write([]byte{b}); err != nil {
			return internal.Error(fmt.Sprintf("write error: %v", x), err, internal.ErrorLoading)
		}
		return nil
	}
	return internal.Error(fmt.Sprintf("not a writer: %v", x), nil, internal.ErrorLoading)
}

type iomap struct {
	fds map[FileDescriptor]filelike
}

func NewIoMap() memory.MemoryAccess {
	fds := map[FileDescriptor]filelike{
		Stdin: filelike{
			descriptor: Stdin,
			access:     Read,
			stream:     os.Stdin,
		},
		Stdout: filelike{
			descriptor: Stdout,
			access:     Write,
			stream:     os.Stdout,
		},
		Stderr: filelike{
			descriptor: Stderr,
			access:     Write,
			stream:     os.Stderr,
		},
	}
	return iomap{fds: fds}
}

func memoryAddressToFileDescriptor(at memory.Address) (FileDescriptor, error) {
	if uint64(at) > uint64(_fileDescriptorsLimt) {
		return 0, internal.Error(
			fmt.Sprintf("unknown address %v", at),
			nil,
			internal.ErrorLoading)
	}
	return FileDescriptor(byte(at)), nil
}

func (x iomap) GetByte(at memory.Address) (byte, error) {
	key, err := memoryAddressToFileDescriptor(at)
	if err != nil {
		return 0, internal.Error(
			fmt.Sprintf("not a descriptor: %v", x),
			err,
			internal.ErrorLoading)
	}
	if reader, ok := x.fds[key]; ok {
		return reader.Read()
	}
	return 0, internal.Error(
		fmt.Sprintf("not a descriptor: %v", x),
		err,
		internal.ErrorLoading)

}
func (x iomap) SetByte(at memory.Address, b byte) error {
	key, err := memoryAddressToFileDescriptor(at)
	if err != nil {
		return internal.Error(
			fmt.Sprintf("not a descriptor: %v", x),
			err,
			internal.ErrorLoading)
	}
	if writer, ok := x.fds[key]; ok {
		return writer.Write(b)
	}
	return internal.Error(
		fmt.Sprintf("not a descriptor: %v", x),
		err,
		internal.ErrorLoading)

}

func (x iomap) GetUint16(memory.Address) (uint16, error) { return 0, nil }
func (x iomap) SetUint16(memory.Address, uint16) error   { return nil }
