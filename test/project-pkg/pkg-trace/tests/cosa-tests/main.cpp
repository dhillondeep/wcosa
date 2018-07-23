#include <uart.h>
#include <trace.h>
#include <Cosa/RTT.hh>

#if !defined(WIO_FRAMEWORK_COSA)
#error "Must be compiled with Cosa"
#endif

void setup() {
    serial::init(9600);
}

void loop() {
    trace::trace << "Hello World!" << trace::endl;
}
