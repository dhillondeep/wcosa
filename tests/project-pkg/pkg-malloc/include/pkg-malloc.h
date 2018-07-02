#ifndef WIO_TESTS_PKGMALLOC_H
#define WIO_TESTS_PKGMALLOC_H

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
