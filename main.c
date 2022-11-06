#include "common.h"
#include "chunk.h"
#include "opcode.h"

int main(int argc, char* argv[]) {
    Chunk chunk;
    initChunk(&chunk);
    writeChunk(&chunk, OP_RETURN);
    freeChunk(&chunk);
}