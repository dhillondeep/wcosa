#include "Arduino.h"
#include <honly/honly.hh>

void setup() {
    Serial.begin(9600);
    Serial.println(sayHonly());
}

void loop() {}

