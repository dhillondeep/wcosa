#ifndef __WIO_TESTS_TRACE_H__
#define __WIO_TESTS_TRACE_H__

namespace serial {
    namespace det {
        struct endl_t {};
        struct trace_t {};
    }
    static det::endl_t endl;
    static det::trace_t trace;

    extern det::trace_t &operator<<(det::trace_t &, const char *str);
    extern det::trace_t &operator<<(det::trace_t &, det::endl_t);
}

#endif
