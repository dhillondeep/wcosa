#include <pkg-list.h>
#include <stdio.h>

int main(void) {
    wio::dynamic_stack stack(16);
    stack.append(5);
    stack.append(10);
    stack.append(15);
    stack.append(20);
    int traverse[] = {20, 15, 10, 5};
    for (int i = 0; i < 4; ++i) {
        if (traverse[i] != stack.pop())
        { return -1; }
    }
    return 0;
}
