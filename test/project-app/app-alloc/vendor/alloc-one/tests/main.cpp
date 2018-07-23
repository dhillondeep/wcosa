#include <alloc-one.h>

int main(void) {
    void *ptr = one::alloc(256);
    void *null = one::alloc(1);
    if (ptr == nullptr) { return 1; }
    if (null != nullptr) { return 1; }
    return 0;
}
