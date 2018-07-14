#include <pkg-fives.h>

void wio::pkg::five_fill(int *arr, int length) {
    constexpr int fill = 5;
    for (int i = 0; i < length; ++i)
    { arr[i] = fill; }
}
