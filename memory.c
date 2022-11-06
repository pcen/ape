#include <stdlib.h>

#include "memory.h"

int growCapacity(int capacity) {
    return capacity < 8 ? : capacity * 2;
}

void* reallocate(void* pointer, size_t oldSize, size_t newSize) {
    if (newSize == 0) {
        free(pointer);
        return NULL;
    }

    void* result = realloc(pointer, newSize);
    if (NULL == result) exit(1);
    return result;
}