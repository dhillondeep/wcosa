#include <stdio.h>
#include <sqlite3.h>
#include <iostream>
#include <thread>

void func1() {
    for (int i = 0; i < 20; i++) {
        std::cout << "Working with DB!!" << std::endl;
    }
}

int main() {
    std::thread thread1(func1);

    sqlite3 *db;
    sqlite3_stmt *stmt;

    sqlite3_open("expenses.db", &db);

    if (db == nullptr) {
        std::cout << "Failed to open DB" << std::endl;
    } else {
        std::cout << "Successfully opened DB" << std::endl;
    }

    thread1.join();
}

