libhonly.so: honly.o
	g++ -shared -o $@ $^

honly.o: honly.cpp
	g++ -c -Wall -fpic $^ -I include -o $@
