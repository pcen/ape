CC=gcc
CFLAGS=-I.
OBJ=main.o chunk.o memory.o debug.o log.o

ape: $(OBJ)
	$(CC) -o ape $(OBJ)

clean:
	rm $(OBJ)
	rm ape 