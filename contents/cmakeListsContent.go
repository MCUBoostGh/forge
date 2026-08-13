package contents

const CMakeListsContent = `
cmake_minimum_required(VERSION 3.20)
project(MyProject C)

set(CMAKE_C_STANDARD 11)

add_executable(MyProject main.c)
`
