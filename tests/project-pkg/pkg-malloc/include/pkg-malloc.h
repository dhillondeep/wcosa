#ifndef WIO_TESTS_PKGMALLOC_H
#define WIO_TESTS_PKGMALLOC_H

#ifndef STACK_SIZE
#error "STACK_SIZE must be defined"
#endif

#if STACK_SIZE == 256
#define memory memory_256
#define stack_alloc stack_alloc_256
#define stack_remaining stack_remaining_256
#define stack_reset stack_reset_256
#elif STACK_SIZE == 512
#define memory memory_512
#define stack_alloc stack_alloc_512
#define stack_remaining stack_remaining_512
#define stack_reset stack_reset_512
#elif STACK_SIZE == 1024
#define memory memory_1024
#define stack_alloc stack_alloc_1024
#define stack_remaining stack_remaining_1024
#define stack_reset stack_reset_1024
#else
#error "STACK_SIZE must be 256, 512, or 1024"
#endif

#ifdef __cplusplus
extern "C" {
#endif
    void *stack_alloc(int size);
    void stack_reset();
    int stack_remaining();
#ifdef __cplusplus
}
#endif

#endif
