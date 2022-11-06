#ifndef ape_memory_h
#define ape_memory_h

// Calculate the new capacity of a dynamic array
int growCapacity(int capacity);

// Reallocate a pointers memory to the size of newSize
void* reallocate(void* pointer, size_t oldSize, size_t newSize);

#endif