gcc -c -Wall -fpic -lm -lyaml -lzmq *.c cJSON/cJSON.c
gcc -shared -o libmango.so *.o
rm *.o
