#include <alloc-two.h>
#include <pkg-malloc.h>

void *two::alloc(int bytes)
{ return stack_alloc(bytes); }
