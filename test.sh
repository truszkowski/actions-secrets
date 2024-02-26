#!/bin/bash

if [ "${TEST_VAR}" = "xxx" ]
then
    echo "OK"
else
    echo "FAIL"
fi

if [ "${TEST_SECRET_TEST}" = "yyy" ]
then
    echo "OK"
else
    echo "FAIL"
fi