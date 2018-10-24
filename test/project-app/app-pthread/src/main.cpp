#include <iostream>
#include <thread>

int foo1 = 0;
int bar2 = 1;

void foo() {
  while (foo1 < 3000) {
    std::cout << "Foo running: " << foo1 << std::endl;
    ++foo1;
  }
}

void bar() {
  while (bar2 < 3000) {
    std::cout << "Bar running: " << bar2 << std::endl;
    ++bar2;
  }
}

int main() {
  std::thread first (foo);     // spawn new thread that calls foo()
  std::thread second (bar);  // spawn new thread that calls bar(0)

  std::cout << "main, foo and bar now execute concurrently...\n";

  // synchronize threads:
  first.join();                // pauses until first finishes
  second.join();               // pauses until second finishes

  std::cout << "foo and bar completed.\n";

  return 0;
}
