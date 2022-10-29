package debug

import (
	"fmt"
	"strconv"
	"the-machine/machine/internal"
	"the-machine/machine/memory"
)

type Action uint8

const (
	Tick        Action = 0
	Next        Action = iota
	Inspect     Action = iota
	PeekRam     Action = iota
	PeekRom     Action = iota
	Registers   Action = iota
	Disassemble Action = iota
	Stack       Action = iota
	Dump        Action = iota
	Load        Action = iota
	Reset       Action = iota
	Quit        Action = iota
)

type Command struct {
	Action Action
}

func (x Command) GetAction() Action {
	return x.Action
}

type Actionable interface {
	GetAction() Action
}

type PeekCommand struct {
	Command
	At     memory.Address
	Length int
}

func NewPeekCommand(action Action, raw string) PeekCommand {
	var pos uint16 = 0
	if at, err := strconv.Atoi(raw); err == nil {
		pos = uint16(at)
	}
	return PeekCommand{
		Command: Command{
			Action: action,
		},
		At:     memory.Address(pos),
		Length: 8,
	}
}

type Interface struct{}

func NewInterface() *Interface {
	return &Interface{}
}

func (x Interface) Prompt(ticks int, ip uint16) {
	fmt.Printf("[tick: %d|ip: %d] > ", ticks, ip)
}

func (x Interface) GetCommand() (Actionable, error) {
	var input string
	fmt.Scanln(&input)
	return x.parseCommand(input)
}

func (x Interface) parseCommand(input string) (Actionable, error) {
	if "" == input {
		return Command{Action: Tick}, nil
	}
	switch input[:1] {
	case "q":
		return Command{Action: Quit}, nil
	case "n":
		return Command{Action: Next}, nil
	case "m":
		if len(input) > 1 {
			return NewPeekCommand(PeekRam, input[1:]), nil
		}
		return Command{Action: PeekRam}, nil
	case "p":
		if len(input) > 1 {
			return NewPeekCommand(PeekRom, input[1:]), nil
		}
		return Command{Action: PeekRom}, nil
	case "s":
		return Command{Action: Stack}, nil
	case "d":
		if len(input) > 3 && input[:4] == "dump" {
			return Command{Action: Dump}, nil
		} else if len(input) > 1 {
			return NewPeekCommand(Disassemble, input[1:]), nil
		} else {
			return Command{Action: Disassemble}, nil
		}
	case "r":
		if len(input) > 4 && input[:5] == "reset" {
			return Command{Action: Reset}, nil
		} else {
			return Command{Action: Registers}, nil
		}
	case "i":
		return Command{Action: Inspect}, nil
	case "l":
		return Command{Action: Load}, nil
	}
	fmt.Printf("Got input: [%s]", input)
	return Command{}, internal.Error(fmt.Sprintf("ERROR: unable to parse command"), nil, internal.ErrorInterface)
}
