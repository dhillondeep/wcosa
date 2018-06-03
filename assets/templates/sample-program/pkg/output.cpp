#include <Cosa/OutputPin.hh>
#include <output.h>

OutputPin ledPin(Board::LED);

void builtLedOn() {
    ledPin.on();
}

void builtLedOff() {
    ledPin.off();
}
