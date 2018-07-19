#include <alloc-one.h>
#include <pkg-malloc.h>

void *one::alloc(int bytes)
{ return stack_alloc(bytes); }
