#include <stdio.h>

#include "vm.h"

int main(int argc, char* argv[]) {

    if (argc != 2) {
        printf("Usage: %s binary\n", argv[0]);
        return -1;
    }

    VM vm;

    FILE* file;
    file = fopen(argv[2], "r");

    if (!file) {
        printf("File not found: %s", argv[2]);
        return -1;
    }



    interpret(&vm, file);
    return 0;
}