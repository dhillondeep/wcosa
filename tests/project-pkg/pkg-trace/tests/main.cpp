#include <trace.h>
#include <Cosa/RTT.hh>

void setup() {
    serial::init(9600);
}

void loop() {
    serial::trace << "Hello world" << std::endl;
}
