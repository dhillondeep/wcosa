#ifndef __WIO_TEST_STACK_H__
#define __WIO_TEST_STACK_H__

namespace wio {
    class dynamic_stack {
    public:
        explicit dynamic_stack(int capacity);
        ~dynamic_stack();

        void append(int val);
        int pop();
        int size() const;
        int capacity() const;

    private:
        int m_size;
        int m_capacity;
        int *m_data;
    };
}

#endif
