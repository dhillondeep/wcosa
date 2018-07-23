#include <pkg-list.h>
#include <Cosa/UART.hh>
#include <Cosa/Trace.hh>

static wio::dynamic_stack stack(16);
static int count = 0;

void setup() {
    uart.begin(9600);
    trace.begin(&uart);
}

void loop() {
    if (stack.size() < stack.capacity()) {
        trace << "Appending " << count << endl;
        stack.append(count++);
    } else {
        trace << "Stack is full" << endl;
    }
    delay(1000);
}
