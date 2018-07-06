#include <uart.h>
#include <Cosa/UART.hh>
#include <stdarg.h>
#include <stdio.h>

#ifndef BUFFER_SIZE
#error "BUFFER_SIZE must be defined"
#endif
static_assert(BUFFER_SIZE == 256, "Expected BUFFER_SIZE == 256");

void uart::init(int baud) {
    uart.begin(baud);
}

void uart::printf(const char *fmt, ...) {
    static char buffer[BUFFER_SIZE];
    va_list args;
    int wrt = 0;

    va_begin(args, fmt);
    wrt = vsnprintf(buffer, fmt, args);
    uart.write(buffer, wrt);
    va_end(args);
}
