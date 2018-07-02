#ifndef STACK_SIZE
#error "STACK_SIZE must be defined"
#endif

#include <pkg-malloc.h>

constexpr int stack_size = STACK_SIZE;
static_assert(stack_size == 256, "Expected STACK_SIZE to be 256");

typedef unsigned char byte;
static byte memory[];
static int ptr = 0;

void *stack_alloc(int size) {
    if (size > stack_remaining())
    { return (void *) 0; }

    void *ptr = (void *) memory[ptr];
    ptr += size;
    return ptr;
}

void stack_reset()
{ ptr = 0; }

int stack_remaining()
{ return stack_size - ptr; }
