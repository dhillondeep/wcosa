#include <lib-complex.h>

int main(){
    coroutine<void>::pull_type source{cooperative};
    error_code ec;
    fail(ec);
    std::cout << ec.value() << '\n';
    
    source();

    boost::thread t{thread};
    t.join();

}

