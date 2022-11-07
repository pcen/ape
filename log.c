#include <stdio.h>

#include <string.h>
#include "common.h"
#include "log.h"

char* levelStr[3] = { "DEBUG", "INFO", "ERROR" };

void LOG(LogLevel level, char* msg) {
    printf("LOG LEVEL: %s - %s\n", levelStr[level], msg);
}
