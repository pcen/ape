#ifndef ape_vm_h
#define ape_vm_h

#include "common.h"

#define STACK_MAX UINT8_MAX
#define CONSTANT_POOL_MAX 64

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
} OpCode;

typedef struct {
    double constants[CONSTANT_POOL_MAX];    // Constant pool (ints only)
    double locals[CONSTANT_POOL_MAX];
    double stack[STACK_MAX];    // Stack of values (everything is double rn)
    int stackPointer;   // Points to first open slot
    uint8_t* ip;    // Points to the next instruction to be run
} VM;

void interpret(VM* vm, FILE* binary);

#endif