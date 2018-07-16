#ifndef __WIO_TESTS_TRACE_H__
#define __WIO_TESTS_TRACE_H__

namespace trace {
    namespace det {
        struct endl_t {};
        struct trace_t {};
    }
    static det::endl_t endl;
    static det::trace_t trace;
}

extern trace::det::trace_t &operator<<(trace::det::trace_t &, const char *str);
extern trace::det::trace_t &operator<<(trace::det::trace_t &, trace::det::endl_t);

#endif
