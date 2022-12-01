#include <stdio.h>

#include "vm.h"

int main(int argc, char* argv[]) {

    if (argc != 2) {
        printf("Usage: %s binary\n", argv[0]);
        return -1;
    }

    VM vm;
    initVM(&vm);
    
    FILE* file = fopen(argv[1], "r");

    if (!file) {
        printf("File not found: %s", argv[1]);
        return -1;
    }

    interpret(&vm, file);
    freeVM(&vm);
    return 0;
}