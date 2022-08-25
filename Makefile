CXX      = g++
CXXFLAGS = -g -Wall -std=c++17

SRC   = ./src
BUILD = ./build

INCLUDES = -I./src

SOURCE_NAMES = token.cc lexer.cc
SOURCES = $(addprefix ./src/, $(SOURCE_NAMES))

OBJECT_NAMES = $(SOURCE_NAMES:.cc=.o)
OBJECTS = $(addprefix ./build/, $(OBJECT_NAMES))

main: $(OBJECTS) main.o
	$(CXX) $(CXXFLAGS) $(OBJECTS) $(BUILD)/main.o -o ape

# place object files in build directory
$(BUILD)/%.o: $(SRC)/%.cc
	$(CXX) $(CXXFLAGS) -c $<  -o $@

# for main
main.o:
	$(CXX) $(CXXFLAGS) -c main.cc  -o $(BUILD)/main.o

.PHONY: clean

clean:
	rm ./build/*
	rm ./ape
