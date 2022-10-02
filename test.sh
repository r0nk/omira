#!/bin/bash

go build .
./omira
./omira add task
./omira add dance
./omira add qwerty
./omira task
echo
./omira status
./omira finish dance
./omira status
./omira finish task
./omira status
./omira finish qwerty
echo
./omira status

#cp ~/life/omira.db .
#./omira schedule
#echo
#./omira status
