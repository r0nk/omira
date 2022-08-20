#!/bin/bash

rm omira.db
go build .
./omira
./omira add task
./omira add dance
./omira add qwerty
./omira task
./omira schedule
echo
./omira status
./omira finish dance
./omira status
./omira finish task
./omira status
./omira finish qwerty
echo
./omira status

