package debug

import "fmt"

type Action uint8

const (
	Tick        Action = 0
	Step        Action = iota
	Inspect     Action = iota
	PeekRam     Action = iota
	PeekRom     Action = iota
	Registers   Action = iota
	Disassemble Action = iota
	Dump        Action = iota
	Quit        Action = iota
)

type Command struct {
	Action Action
}

type Interface struct{}

func NewInterface() *Interface {
	return &Interface{}
}

func (x Interface) Prompt(ticks int) {
	fmt.Printf("[tick: %d] > ", ticks)
}

func (x Interface) GetCommand() (Command, error) {
	var input string
	fmt.Scanln(&input)
	return x.parseCommand(input)
}

func (x Interface) parseCommand(input string) (Command, error) {
	if "" == input {
		return Command{Action: Tick}, nil
	}
	switch input[:1] {
	case "q":
		return Command{Action: Quit}, nil
	case "s":
		return Command{Action: Step}, nil
	case "m":
		return Command{Action: PeekRam}, nil
	case "p":
		return Command{Action: PeekRom}, nil
	case "d":
		if len(input) > 3 && input[:4] == "dump" {
			return Command{Action: Dump}, nil
		} else {
			return Command{Action: Disassemble}, nil
		}
	case "r":
		return Command{Action: Registers}, nil
	case "i":
		return Command{Action: Inspect}, nil
	}
	fmt.Printf("Got input: [%s]", input)
	return Command{}, fmt.Errorf("ERROR: unable to parse command")
}
