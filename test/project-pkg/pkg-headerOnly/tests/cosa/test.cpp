#include <Cosa/UART.hh>
#include <Cosa/Trace.hh>
#include <honly/honly.hh>

void setup() {
    uart.begin(9600);
    trace.begin(&uart);

    trace << sayHonly() << endl;
}

void loop() {}

