#include <uart.h>
#include <Cosa/RTT.hh>

void setup()
{ serial::init(9600); }

void loop() {
    delay(500);
    serial::printf("Hello %s\n", "world");
}
