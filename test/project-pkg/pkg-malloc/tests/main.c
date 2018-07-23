#include <pkg-malloc.h>

#ifndef NULL
#define NULL ((void *) 0)
#endif

int main(void) {
    {
        void *ptr = stack_alloc(128);
        if (ptr == NULL)
        { return 1; }

        int expected = 128;
        int remaining = stack_remaining();
        if (remaining != expected)
        { return 1; }

        stack_reset();
        expected = 256;
        remaining = stack_remaining();
        if (remaining != expected)
        { return 1; }
    }
    {
        void *ptr = stack_alloc(512);
        if (ptr != NULL)
        { return 1; }

        int expected = 256;
        int remaining = stack_remaining();
        if (remaining != expected)
        { return 1; }

        stack_reset();
    }
    {
        int *ptr = (int *) stack_alloc(sizeof(int));
        *ptr = 13;
        if (*ptr != 13)
        { return 1; }

        *ptr = 25;
        if (*ptr != 25)
        { return 1; }

        if (stack_remaining() != 256 - sizeof(int))
        { return 1; }

        stack_reset();
        if (stack_remaining() != 256)
        { return 1; }
    }
    return 0;
}
