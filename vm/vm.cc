#include "vm.h"

#include <fstream>
#include <iostream>

VirtualMachine::VirtualMachine()
	: pc(0), sp(0) {}

VirtualMachine::~VirtualMachine() {}

i32 readInt32(std::ifstream& ifs) {
	i32 temp;
	ifs.read(reinterpret_cast<char*>(&temp), sizeof(i32));
	return temp;
}

void VirtualMachine::load(const char* file) {
	std::ifstream ifs(file, std::ios::in | std::ios::binary);
	ifs.unsetf(std::ios::skipws); // ifstream still eats newlines in binary mode

	ifs.seekg(0, std::ios::end);
	std::streampos size = ifs.tellg();
	ifs.seekg(0, std::ios::beg);

	i32 numLits = readInt32(ifs);
	std::cout << "num literals: " << numLits << std::endl;

	for (i32 i = 0; i < numLits; i++) {
		literals.push_back(readInt32(ifs));
	}

	this->bytecode.reserve(size); // does not account for data portions of file
	std::copy(std::istream_iterator<u8>(ifs), std::istream_iterator<u8>(), std::back_inserter(bytecode));
	for (auto c : bytecode) {
		std::cout << "file: op: " << (int)c << std::endl;
	}

	ifs.close();
}

void VirtualMachine::interpret(const char* file) {
	load(file);
	run();
}

i32 VirtualMachine::pop() {
	i32 value = stack.top();
	stack.pop();
	return value;
}

void VirtualMachine::push(i32 value) {
	stack.emplace(value);
}

u8 VirtualMachine::next() {
	return bytecode[pc++];
}

void VirtualMachine::run() {
	u8 op = next();
	while(op != 0) {
		std::cerr << "run: op: " << (int) op << std::endl;
		switch (op) {
		case OP_NIL:
			break;
		case OP_SET: {
			u8 idx = next();
			locals[idx] = pop();
			break;
		}
		case OP_GET: {
			u8 idx = next();
			push(locals[idx]);
			break;
		}
		case OP_ADD: {
			i32 b = pop();
			i32 a = pop();
			push(a + b);
			break;
		}
		case OP_SUBTRACT: {
			i32 b = pop();
			i32 a = pop();
			push(a - b);
			break;
		}
		case OP_MULTIPLY: {
			i32 b = pop();
			i32 a = pop();
			push(a * b);
			break;
		}
		case OP_DIVIDE: {
			i32 b = pop();
			i32 a = pop();
			push(a / b);
			break;
		}
		case OP_CONSTANT: {
			u8 idx = next();
			push(literals[idx]);
			break;
		}
		case OP_PRINT:
			std::cout << pop() << std::endl;
			break;
		default:
			std::cout << "error: unknown opcode: " << (int)op << std::endl;
			exit(1);
		}
		op = next();
	}
}
