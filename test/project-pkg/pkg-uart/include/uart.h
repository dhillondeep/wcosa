#ifndef __WIO_TESTS_UART_H__
#define __WIO_TESTS_UART_H__

#ifndef __cplusplus
#error "Buildable only with C++"
#endif

namespace serial {
    extern void init(int baud);
    extern void printf(const char *fmt, ...);
}

#endif
