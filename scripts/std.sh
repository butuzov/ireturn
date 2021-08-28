#!/usr/bin/env bash

mkdir -p .tmp

# create std pkg list
for i in $(seq 1 17); do
    docker run --rm -it golang:1.$i go list std \
        | grep -v internal | grep -v vendor > .tmp/go@1.$i;
    # docker rmi golang:1.$i
done


echo "// Code generated using std.sh; DO NOT EDIT." > "std.go"
echo "" >> "std.go"

echo "// We will ignore that fact that some of packages" >> "std.go"
echo "// were removed from stdlib." >> "std.go"
echo "" >> "std.go"
echo "package ireturn" >> "std.go"
echo "" >> "std.go"
echo "" >> "std.go"
echo "var std = []string{" >> "std.go"

for i in $(seq 2 17); do

    printf "\t// added in Go v1.%s in compare to v1.%s (docker image)\n" $i "$(($i-1))" >> "std.go"

    for pkg in $(comm -13 <(sort .tmp/go@1.$(($i-1))) <(sort .tmp/go@1.$i)); do
        printf "\t\"%s\",\n" $(echo $pkg | tr -d \\r) >> "std.go"
    done

    # whats removed in new release not really required.
    # printf "(removed) Go v1.%s vs v1.%s\n" $i "$(($i-1))"
    # comm -13  <(sort go@1.$i) <(sort go@1.$(($i-1)))
done

echo "}" >> "std.go"
