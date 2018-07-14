#include <uart.h>
#include <Cosa/RTT.hh>

static int count = 0;

void setup() {
    RTT::begin();
}

void loop() {
    serial::printf("Loop number: %d\n", count);
}
