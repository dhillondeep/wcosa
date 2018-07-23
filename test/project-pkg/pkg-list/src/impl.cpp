#include <pkg-list.h>
#include <pkg-malloc.h>

using namespace wio;

dynamic_stack::dynamic_stack(int capacity) :
    m_size(0),
    m_capacity(capacity),
    m_data(static_cast<int *>(stack_alloc(capacity))) {}

dynamic_stack::~dynamic_stack()
{ stack_reset(); }

void dynamic_stack::append(int val) {
    if (m_size >= m_capacity) { return; }
    m_data[m_size++] = val;
}

int dynamic_stack::pop() {
    if (m_size == 0) { return 0; }
    --m_size;
    return m_data[m_size];
}

int dynamic_stack::size() const
{ return m_size; }

int dynamic_stack::capacity() const
{ return m_capacity; }
