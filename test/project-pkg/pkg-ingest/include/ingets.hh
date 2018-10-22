#include <iostream>

using namespace std;

void display() {
    cout << "Definition \"Req\"        =" << Req << endl;
    cout << "Definition \"Ingest\"     =" << Ingest << endl;

#ifdef Opt
    cout << "Definition \"Opt\"        =" << Opt << endl;
#endif
}

