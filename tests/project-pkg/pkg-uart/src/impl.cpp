#include <uart.h>
#include <Cosa/UART.hh>
#include <stdarg.h>
#include <stdio.h>

#ifndef BUFFER_SIZE
#error "BUFFER_SIZE must be defined"
#endif
static_assert(BUFFER_SIZE == 256, "Expected BUFFER_SIZE == 256");

void serial::init(int baud) {
    uart.begin(baud);
}

void serial::printf(const char *fmt, ...) {
    static char buffer[BUFFER_SIZE];
    va_list args;
    int wrt = 0;

    va_start(args, fmt);
    wrt = vsprintf(buffer, fmt, args);
    uart.write(buffer, wrt);
    va_end(args);
}
