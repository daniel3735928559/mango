gcc -g -c -Wall -fpic -lm -lyaml -lzmq *.c cJSON/cJSON.c
x=$?
gcc -g -shared -o libmango.so *.o
rm *.o
if [[ $x -eq 0 ]]; then
echo DONE
else
echo ERROR
fi
