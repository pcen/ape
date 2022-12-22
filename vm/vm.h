#pragma once

#include "common.h"

#include <vector>
#include <stack>
#include <array>
#include <cstdint>

typedef enum {
	OP_NIL,
	OP_SET,
	OP_GET,
	OP_ADD,
	OP_SUBTRACT,
	OP_MULTIPLY,
	OP_DIVIDE,
	OP_CONSTANT,
	OP_PRINT
} Opcode;

class VirtualMachine {
public:
	VirtualMachine();
	~VirtualMachine();

	void load(const char* file);
	void interpret(const char* file);

private:
	void loadBytecode();
	void run();

	i32 pop(); // pop from stack
	void push(i32); // push to stack
	u8 next(); // get next opcode

	std::vector<i32> literals;
	std::array<i32, 512> locals;
	std::stack<i32> stack;

	int sp; // stack pointer: first open slot of stack
	int pc; // program counter
	std::vector<u8> bytecode;
};
