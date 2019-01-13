#include "show.hh"
#include <ingest.hh>
#include <iostream>


void show() {
    std::cout << "User Name: " << STRINGIFYMACRO(USER_NAME) << std::endl;
    std::cout << "User City: " << STRINGIFYMACRO(USER_CITY) << std::endl;

#ifdef USER_AGE
    std::cout << "User Age: " << STRINGIFYMACRO(USER_AGE) << std::endl;
#endif
#ifdef LIVES_IN_CANADA
    std::cout << "User lives in Canada" << std::endl;
#else
    std::cout << "User does not live in Canada" << std::endl;
#endif

    std::cout << "\nDisplay from Ingest package: " << std::endl;
    display();
}

