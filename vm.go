package main

import (
	"fmt"
	"encoding/binary"
)

// Opcode represents the operation to be performed
type Opcode byte

const (
	PUSH Opcode = iota
	POP
	ADD
	SUB
	MUL
	DIV
	JMP
	JMPZ
	HALT
)

// Instruction represents a single VM instruction
type Instruction struct {
	Opcode Opcode
	Operand int64
}

// VM represents our simple virtual machine
type VM struct {
	instructions []Instruction
	stack        []int64
	pc           int
}

// NewVM creates a new VM instance
func NewVM(instructions []Instruction) *VM {
	return &VM{
		instructions: instructions,
		stack:        make([]int64, 0),
		pc:           0,
	}
}

// Run executes the VM instructions
func (vm *VM) Run() {
	for vm.pc < len(vm.instructions) {
		inst := vm.instructions[vm.pc]
		switch inst.Opcode {
		case PUSH:
			vm.stack = append(vm.stack, inst.Operand)
		case POP:
			vm.stack = vm.stack[:len(vm.stack)-1]
		case ADD:
			b, a := vm.stack[len(vm.stack)-1], vm.stack[len(vm.stack)-2]
			vm.stack = vm.stack[:len(vm.stack)-2]
			vm.stack = append(vm.stack, a+b)
		case SUB:
			b, a := vm.stack[len(vm.stack)-1], vm.stack[len(vm.stack)-2]
			vm.stack = vm.stack[:len(vm.stack)-2]
			vm.stack = append(vm.stack, a-b)
		case MUL:
			b, a := vm.stack[len(vm.stack)-1], vm.stack[len(vm.stack)-2]
			vm.stack = vm.stack[:len(vm.stack)-2]
			vm.stack = append(vm.stack, a*b)
		case DIV:
			b, a := vm.stack[len(vm.stack)-1], vm.stack[len(vm.stack)-2]
			vm.stack = vm.stack[:len(vm.stack)-2]
			vm.stack = append(vm.stack, a/b)
		case JMP:
			vm.pc = int(inst.Operand) - 1
		case JMPZ:
			if vm.stack[len(vm.stack)-1] == 0 {
				vm.pc = int(inst.Operand) - 1
			}
			vm.stack = vm.stack[:len(vm.stack)-1]
		case HALT:
			return
		}
		vm.pc++
	}
}

// EBPFInstruction represents an eBPF instruction
type EBPFInstruction struct {
	Opcode byte
	Regs   byte
	Offset int16
	Imm    int32
}

// CompileToEBPF compiles VM instructions to eBPF bytecode
func CompileToEBPF(instructions []Instruction) []EBPFInstruction {
	var ebpfInstructions []EBPFInstruction
	for _, inst := range instructions {
		switch inst.Opcode {
		case PUSH:
			ebpfInstructions = append(ebpfInstructions, EBPFInstruction{
				Opcode: 0xb7, // mov64 r0, imm
				Regs:   0x00,
				Imm:    int32(inst.Operand),
			})
			ebpfInstructions = append(ebpfInstructions, EBPFInstruction{
				Opcode: 0x7b, // stxdw [r10-8], r0
				Regs:   0xa0,
				Offset: -8,
			})
		case POP:
			ebpfInstructions = append(ebpfInstructions, EBPFInstruction{
				Opcode: 0x79, // ldxdw r0, [r10-8]
				Regs:   0xa0,
				Offset: -8,
			})
		case ADD:
			ebpfInstructions = append(ebpfInstructions, EBPFInstruction{
				Opcode: 0x79, // ldxdw r0, [r10-8]
				Regs:   0xa0,
				Offset: -8,
			})
			ebpfInstructions = append(ebpfInstructions, EBPFInstruction{
				Opcode: 0x79, // ldxdw r1, [r10-16]
				Regs:   0xa1,
				Offset: -16,
			})
			ebpfInstructions = append(ebpfInstructions, EBPFInstruction{
				Opcode: 0x0f, // add64 r0, r1
				Regs:   0x10,
			})
			ebpfInstructions = append(ebpfInstructions, EBPFInstruction{
				Opcode: 0x7b, // stxdw [r10-16], r0
				Regs:   0xa0,
				Offset: -16,
			})
		// Add similar cases for SUB, MUL, DIV
		case JMP:
			ebpfInstructions = append(ebpfInstructions, EBPFInstruction{
				Opcode: 0x05, // ja +offset
				Offset: int16(inst.Operand),
			})
		case JMPZ:
			ebpfInstructions = append(ebpfInstructions, EBPFInstruction{
				Opcode: 0x79, // ldxdw r0, [r10-8]
				Regs:   0xa0,
				Offset: -8,
			})
			ebpfInstructions = append(ebpfInstructions, EBPFInstruction{
				Opcode: 0x15, // jeq r0, 0, +offset
				Regs:   0x00,
				Offset: int16(inst.Operand),
			})
		case HALT:
			ebpfInstructions = append(ebpfInstructions, EBPFInstruction{
				Opcode: 0x95, // exit
				Regs:   0x00,
			})
		}
	}
	return ebpfInstructions
}

func vm() {
	// Example program: calculate (10 + 5) * 2
	program := []Instruction{
		{PUSH, 10},
		{PUSH, 5},
		{ADD, 0},
		{PUSH, 2},
		{MUL, 0},
		{HALT, 0},
	}

	// Run the program in our VM
	vm := NewVM(program)
	vm.Run()
	fmt.Printf("VM Result: %d\n", vm.stack[len(vm.stack)-1])

	// Compile to eBPF
	ebpfInstructions := CompileToEBPF(program)

	// Print eBPF bytecode
	fmt.Println("eBPF Bytecode:")
	for _, inst := range ebpfInstructions {
		bytes := make([]byte, 8)
		binary.LittleEndian.PutUint64(bytes, uint64(inst.Opcode)|
			uint64(inst.Regs)<<8|
			uint64(uint16(inst.Offset))<<16|
			uint64(uint32(inst.Imm))<<32)
		fmt.Printf("%02x %02x %02x %02x %02x %02x %02x %02x\n",
			bytes[0], bytes[1], bytes[2], bytes[3],
			bytes[4], bytes[5], bytes[6], bytes[7])
	}
}