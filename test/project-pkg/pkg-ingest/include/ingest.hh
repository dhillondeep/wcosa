#include <iostream>

#define STRINGIFY(x) #x
#define STRINGIFYMACRO(y) STRINGIFY(y)

using namespace std;

void display() {
    cout << "Definition \"Req\"        =" << STRINGIFYMACRO(Req) << endl;
    cout << "Definition \"Ingest\"     =" << STRINGIFYMACRO(Ingest) << endl;

#ifdef Opt
    cout << "Definition \"Opt\"        =" << STRINGIFYMACRO(Opt) << endl;
#endif

#ifdef LIVES_IN_CANADA
    cout << "Package called and User lives in Canada" << endl;
#endif

#ifdef USER_AGE
    cout << "Package called and User age is " << STRINGIFYMACRO(USER_AGE) << endl;
#endif
}

