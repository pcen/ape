#include <stdlib.h>

#include "chunk.h"
#include "memory.h"

void initChunk(Chunk* chunk) {
    chunk->size = 0;
    chunk->capacity = 0;
    chunk->code = NULL;
}

void writeChunk(Chunk* chunk, BYTE byte) {
    if (chunk -> size == chunk->capacity) {
        int newCapacity = growCapacity(chunk->capacity);
        chunk->code = (BYTE*) reallocate(chunk->code, chunk->capacity, newCapacity);
        chunk->capacity = newCapacity;
    }

    chunk->code[chunk->size] = byte;
    chunk->size++;
}

void freeChunk(Chunk* chunk) {
    reallocate(chunk->code, sizeof(BYTE) * chunk->capacity, 0);
    initChunk(chunk);
}