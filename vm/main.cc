#include <stdio.h>

#include "vm.h"

int main(int argc, char* argv[]) {

	if (argc != 2) {
		printf("Usage: %s binary\n", argv[0]);
		return -1;
	}

	VirtualMachine vm;
	vm.interpret(argv[1]);

	return 0;
}
