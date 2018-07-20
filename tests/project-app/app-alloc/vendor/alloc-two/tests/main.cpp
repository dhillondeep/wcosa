#include <alloc-two.h>

int main(void) {
    void *ptr = two::alloc(256);
    void *null = two::alloc(1);
    if (ptr == nullptr) { return 1; }
    if (null != nullptr) { return 1; }
    return 0;
}
