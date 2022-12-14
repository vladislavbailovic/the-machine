package device

import (
	"encoding/binary"
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

	_fileDescriptorsLimt FileDescriptor = 255
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

type Filelike struct {
	descriptor FileDescriptor
	access     AccessType
	stream     interface{}
}

func NewFilelike(fd FileDescriptor, access AccessType, stream interface{}) Filelike {
	return Filelike{
		descriptor: fd,
		access:     access,
		stream:     stream,
	}
}

func (x Filelike) String() string {
	stream := "<stream>"
	if x.stream == nil {
		stream = "<NO STREAM>"
	}
	return fmt.Sprintf("%s <%s>: %s", x.descriptor, x.access, stream)
}

func (x Filelike) Read() (byte, error) {
	if x.access != Read {
		return 0, internal.Error(fmt.Sprintf("unable to read file descriptor in %v", x), nil, internal.ErrorLoading)
	}
	if reader, ok := x.stream.(io.Reader); ok {
		buf := make([]byte, 1, 1)
		if n, err := reader.Read(buf); err != nil {
			if io.EOF == err {
				return 0, nil
			}
			return 0, internal.Error(fmt.Sprintf("read error: %v", x), err, internal.ErrorLoading)
		} else if n > 0 {
			return buf[0], nil
		}
		return 0, internal.Error(fmt.Sprintf("read error: EOF %v", x), nil, internal.ErrorLoading)
	}
	return 0, internal.Error(fmt.Sprintf("not a reader: %v", x), nil, internal.ErrorLoading)
}

func (x Filelike) Write(b byte) error {
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

func (x Filelike) WriteUint16(b []byte) error {
	if x.access != Write {
		return internal.Error(fmt.Sprintf("unable to write to file descriptor in %v", x), nil, internal.ErrorLoading)
	}
	if writer, ok := x.stream.(io.Writer); ok {
		if _, err := writer.Write(b); err != nil {
			return internal.Error(fmt.Sprintf("write error: %v", x), err, internal.ErrorLoading)
		}
		return nil
	}
	return internal.Error(fmt.Sprintf("not a writer: %v", x), nil, internal.ErrorLoading)
}

type IOMap struct {
	fds map[FileDescriptor]Filelike
}

func NewIoMap() memory.MemoryAccess {
	fds := map[FileDescriptor]Filelike{
		Stdin: Filelike{
			descriptor: Stdin,
			access:     Read,
			stream:     os.Stdin,
		},
		Stdout: Filelike{
			descriptor: Stdout,
			access:     Write,
			stream:     os.Stdout,
		},
		Stderr: Filelike{
			descriptor: Stderr,
			access:     Write,
			stream:     os.Stderr,
		},
	}
	return &IOMap{fds: fds}
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

func (x IOMap) GetByte(at memory.Address) (byte, error) {
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
func (x IOMap) SetByte(at memory.Address, b byte) error {
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

func (x IOMap) GetUint16(at memory.Address) (uint16, error) {
	val, err := x.GetByte(at)
	if err != nil {
		return uint16(val), err
	}
	return uint16(val), nil
}
func (x IOMap) SetUint16(at memory.Address, what uint16) error {
	key, err := memoryAddressToFileDescriptor(at)
	if err != nil {
		return internal.Error(
			fmt.Sprintf("not a descriptor: %v", x),
			err,
			internal.ErrorLoading)
	}

	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, what)

	if writer, ok := x.fds[key]; ok {
		return writer.WriteUint16(b)
	}
	return internal.Error(
		fmt.Sprintf("not a descriptor: %v", x),
		err,
		internal.ErrorLoading)
}

func (x *IOMap) SetDescriptor(fd FileDescriptor, what Filelike) {
	x.fds[fd] = what
}
