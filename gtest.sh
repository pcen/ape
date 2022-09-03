#!/bin/bash
GTEST_SOURCE="external/googletest"

mkdir "./build/googletest"
cd "$_"
cmake "../../${GTEST_SOURCE}"
make
