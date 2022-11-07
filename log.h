#ifndef ape_log_h
#define ape_log_h

typedef enum {
    DEBUG,
    INFO,
    ERROR,
} LogLevel;

void LOG(LogLevel, char* msg);

#endif