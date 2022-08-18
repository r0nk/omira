#!/bin/bash

rm omira.db
go build .
./omira
./omira add task
./omira add dance
./omira task
./omira schedule
echo
./omira status
./omira finish dance
./omira finish task
echo
./omira status

