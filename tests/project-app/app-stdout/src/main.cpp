#ifndef OUTPUT_VAL
#error "OUTPUT_VAL must be defined"
#endif

#include <stdio.h>
#include <stdlib.h>

static constexpr int val = OUTPUT_VAL;
static_assert(val > 77, "Expected OUTPUT_VAL > 77");

int main(int argc, char *argv[]) {
    if (argc != 2)
    { return 1; }

    constexpr int base = 10;
    const char *str_count = argv[1];
    char *end = NULL;
    long count = strtol(str_count, &end, 10);
    if (*end != '\0')
    { return 1; }

    for (long i = 0; i < count; ++i)
    { fprintf(stdout, "%i", val); }

    return 0;
}
