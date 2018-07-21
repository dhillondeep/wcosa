#include <stdio.h>
#include <stdint.h>
#include <alloc-one.h>
#include <alloc-two.h>

namespace det {
    template<int ptr_size> struct ptr_type;

    template<> struct ptr_type<1> { typedef uint8_t  type; };
    template<> struct ptr_type<2> { typedef uint16_t type; };
    template<> struct ptr_type<4> { typedef uint32_t type; };
    template<> struct ptr_type<8> { typedef uint64_t type; };
}

struct uintptr {
    typedef typename det::ptr_type<sizeof(void *)>::type type;
};

int main(void) {
    static void *one_ptrs[32] = { nullptr };
    static void *two_ptrs[32] = { nullptr };
    for (int i = 0; i < 32; ++i) {
        one_ptrs[i] = one::alloc(32);
    }
    for (int i = 0; i < 32; ++i) {
        two_ptrs[i] = two::alloc(32);
    }
    printf("one_ptrs = {");
    for (int i = 0; i < 32; ++i) {
        if (i % 4 == 0) { printf("\n\t"); }
        printf("%#08x, ", static_cast<unsigned int>(reinterpret_cast<uintptr::type>(one_ptrs[i])));
    }
    printf("\n}\n\n");
    printf("two_ptrs = {");
    for (int i = 0; i < 32; ++i) {
        if (i % 4 == 0) { printf("\n\t"); }
        printf("%#08x, ", static_cast<unsigned int>(reinterpret_cast<uintptr::type>(two_ptrs[i])));
    }
    printf("\n}\n");

    int p = 0;
    for (; p < 8; ++p) {
        if (one_ptrs[p] == nullptr) { return 1; }
    }
    for (; p < 32; ++p) {
        if (one_ptrs[p] != nullptr) { return 1; }
    }
    for (p = 0; p < 16; ++p) {
        if (two_ptrs[p] == nullptr) { return 1; }
    }
    for (; p < 32; ++p) {
        if (one_ptrs[p] != nullptr) { return 1; }
    }
    return 0;
}
