#include <boost/system/error_code.hpp>
#include <boost/thread.hpp>
#include <boost/chrono.hpp>
#include <iostream>

using namespace boost::system;

void fail(error_code &ec) {
      ec = errc::make_error_code(errc::not_supported);
}

void wait(int seconds){
    boost::this_thread::sleep_for(boost::chrono::seconds{seconds});
}

void thread() {
    for (int i = 0; i < 5; ++i) {
        wait(1);
        std::cout << i << '\n';            
    }
}

