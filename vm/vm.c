#include <stdio.h>
#include <stdlib.h>

#include "vm.h"

// Set the IP of the VM and sets topStack pointer correctly.
static void loadBytecode(VM* vm, FILE* binary) {

    // Get size of file
    fseek(binary, 0L, SEEK_END); // Seek to end
    size_t size = ftell(binary); // Get offset
    fseek(binary, 0L, SEEK_SET); // Back to start

    uint32_t numConstants;
    fread(&numConstants, 4, 1, binary);
    fseek(binary, 4, SEEK_SET); // Move 1 byte forward

    uint32_t constant[numConstants];
    fread(&constant, 4, numConstants, binary);
    for (uint32_t i = 0; i < numConstants; i++) {
        vm->constants[i] = constant[i];
    }
    fseek(binary, 4 * numConstants, SEEK_CUR);  // Now we're at opcode

    // Dump the contents here
    size_t opcodeSize = size - 4 - (numConstants * 4);

    vm->ip = malloc(opcodeSize * sizeof(uint8_t));
    fread(vm->ip, 1, opcodeSize, binary);

    fclose(binary);
}

static double pop(VM* vm) {
    return vm->stack[--vm->stackPointer];
}

static void push(VM* vm, double value) {
    vm->stack[vm->stackPointer++] = value;
}

static uint8_t readByte(VM* vm) {
    return *(vm->ip++);
}

static void run(VM* vm) {
    while(true) {

        uint8_t opcode = *(vm->ip++);

        switch (opcode) {
        case OP_NIL:
            break;
        case OP_SET: {
            uint8_t idx = readByte(vm);
            vm->locals[idx] = pop(vm);
            break;
        }
        case OP_GET: {
            uint8_t idx = readByte(vm);
            push(vm, vm->locals[idx]);
            break;
        }
        case OP_ADD: {
            double b = pop(vm);
            double a = pop(vm);
            push(vm, a + b);
            break;
        }
        case OP_SUBTRACT: {
            double b = pop(vm);
            double a = pop(vm);
            push(vm, a - b);
            break;
        }
        case OP_MULTIPLY: {
            double b = pop(vm);
            double a = pop(vm);
            push(vm, a * b);
            break;
        }
        case OP_DIVIDE: {
            double b = pop(vm);
            double a = pop(vm);
            push(vm, a / b);
            break;
        }
        case OP_CONSTANT: {
            uint8_t idx = readByte(vm);
            push(vm, vm->constants[idx]);
            break;
        }
        case OP_PRINT:
            printf("%f\n", pop(vm));
            break;
        default:
            exit(1);
        }
    }
}

void interpret(VM* vm, FILE* binary) {
    loadBytecode(vm, binary);
    run(vm);
}

