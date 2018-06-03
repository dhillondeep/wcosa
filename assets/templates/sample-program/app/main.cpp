#include <Cosa/RTT.hh>
#include <Cosa/OutputPin.hh>

OutputPin ledPin(Board::LED);

void setup() {
    RTT::begin();
}

void loop() {
    ledPin.on();
    delay(50);
    ledPin.off();
    delay(500);
}
