#ifndef PRINTER_H
#define PRINTER_H

#include <Cosa/Trace.hh>

// Prints Hello and the word specified by the user
// Ex: Hello World
void sayHello(char* word) {
    trace << "Hello " << word << endl;
}

#endif /* PRINTER_H */
