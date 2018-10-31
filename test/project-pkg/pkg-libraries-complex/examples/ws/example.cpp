#include <lib-complex.h>

int main(){
    error_code ec;
    fail(ec);
    std::cout << ec.value() << '\n';
    
    boost::thread t{thread};
    t.join();
}

