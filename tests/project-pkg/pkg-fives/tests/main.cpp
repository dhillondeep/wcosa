#include <pkg-fives.h>

int main(void) {
    {
        constexpr int length = 25;
        static int arr[length];
        int expected = 5;
        wio::pkg::five_fill(arr, length);
        for (int i = 0; i < length; ++i) {
            if (arr[i] != expected)
            { return 1; }
        }
    }
    return 0;
}
