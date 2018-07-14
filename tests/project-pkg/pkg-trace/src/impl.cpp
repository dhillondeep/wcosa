#include <trace.h>
#include <uart.h>

using namespace serial;

det::trace_t &operator<<(det::trace_t &, const char *str) {
    serial::printf(str);
    return serial::trace;
}

det::trace_t &operator<<(det::trace_t &, det::endl_t) {
    serial::printf("\n");
    return serial::trace;
}
