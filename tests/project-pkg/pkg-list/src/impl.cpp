#include <pkg-list.h>
#include <pkg-malloc.h>

using namespace wio;

dynamic_stack::dynamic_stack(int size) :
    m_size(0),
    m_data(static_cast<int *>(stack_alloc(size))) {}

dynamic_stack::~dynamic_stack()
{ stack_reset(); }

void dynamic_stack::append(int val)
{ m_data[m_size++] = val; }

int dynamic_stack::pop() {
    --m_size;
    return m_data[m_size];
}
