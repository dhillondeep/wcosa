#include <Cosa/Trace.hh>
#include <Cosa/UART.hh>
#include <printer.h>

void setup() {
    uart.begin(9600);
    trace.begin(&uart);
}

void loop() {
    sayHello("World");
}

