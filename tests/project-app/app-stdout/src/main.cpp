#include <stdio.h>
#include <stdlib.h>

int main(int argc, char *argv[]) {
    if (argc != 3)
    { return 1; }

    constexpr int base = 10;
    const char *str_count = argv[1];
    const char *str_val = argv[2];
    char *end = NULL;
    long count = strtol(str_count, &end, 10);
    if (*end != '\0')
    { return 1; }
    long val = strtol(str_val, &end, 10);
    if (*end != '\0')
    { return 1; }

    for (long i = 0; i < count; ++i)
    { fprintf(stdout, "%li", val); }
    printf("\n");

    return 0;
}
