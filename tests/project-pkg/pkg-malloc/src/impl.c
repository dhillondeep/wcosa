#ifndef STACK_SIZE
#error "STACK_SIZE must be defined"
#endif

#define static_assert(cond, msg) typedef char __static_assertion[(cond) ? 1 : -1]

#include <pkg-malloc.h>

static_assert(STACK_SIZE >= 256, "Expected STACK_SIZE to be at least 256");

typedef unsigned char byte;
static byte memory[STACK_SIZE];
static int ptr = 0;

void *stack_alloc(int size) {
    if (size > stack_remaining())
    { return (void *) 0; }

    void *mem = (void *) (memory + ptr);
    ptr += size;
    return mem;
}

void stack_reset()
{ ptr = 0; }

int stack_remaining()
{ return STACK_SIZE - ptr; }
