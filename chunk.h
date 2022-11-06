#ifndef ape_chunk_h
#define ape_chunk_h

#include "opcode.h"
#include "common.h"

typedef struct {
    int size;
    int capacity;
    BYTE* code;
} Chunk;


void initChunk(Chunk* chunk);
void writeChunk(Chunk* chunk, BYTE byte);
void freeChunk(Chunk* chunk);

#endif