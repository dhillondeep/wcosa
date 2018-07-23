#include <trace.h>
#include <uart.h>

trace::det::trace_t &operator<<(trace::det::trace_t &, const char *str) {
    serial::printf(str);
    return trace::trace;
}

trace::det::trace_t &operator<<(trace::det::trace_t &, trace::det::endl_t) {
    serial::printf("\n");
    return trace::trace;
}

