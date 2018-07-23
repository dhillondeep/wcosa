#include <pkg-square.h>

int main(void) {
    {
        int input = 5;
        int expected = 25;
        int output = wio::pkg::square(input);
        if (output != expected)
        { return 1; }
    }
    {
        long input = 6;
        long expected = 36;
        long output = wio::pkg::square(input);
        if (output != expected)
        { return 1; }
    }
    return 0;
}
