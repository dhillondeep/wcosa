#include <Cosa/RTT.hh>
#include <output.h>

void setup() {
    RTT::begin();
}

void loop() {
    builtLedOn();
    delay(50);
    builtLedOff();
    delay(500);
}
