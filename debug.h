#ifndef ape_debug_h
#define ape_debug_h

#include "chunk.h"
#include "opcode.h"

void disassembleChunk(Chunk* chunk, const char* name);
int disassembleInstruction(Chunk* chunk, int offset);

#endif