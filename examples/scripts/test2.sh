#!/bin/bash

echo "test"
for arg in "$@"
do
    echo "arg: $arg"
done
echo '----------------""-----'
echo $#
echo $@
exit 0
