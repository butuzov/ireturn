#!/usr/bin/env bash

mkdir -p .tmp

# create std pkg list
for i in $(seq 1 17); do
    docker run --rm -it golang:1.$i go list std \
        | grep -v internal | grep -v vendor > .tmp/go@1.$i;
    # docker rmi golang:1.$i
done

dest="pkg/analyzer/std.go"

echo "// Code generated using std.sh; DO NOT EDIT." > $dest
echo "" >> $dest

echo "// We will ignore that fact that some of packages" >> $dest
echo "// were removed from stdlib." >> $dest
echo "" >> $dest
echo "package analyzer" >> $dest
echo "" >> $dest
echo "var std = map[string]struct{}{" >> $dest

for i in $(seq 2 17); do

    printf "\t// added in Go v1.%s in compare to v1.%s (docker image)\n" $i "$(($i-1))" >> $dest

    for pkg in $(comm -13 <(sort .tmp/go@1.$(($i-1))) <(sort .tmp/go@1.$i)); do
        printf "\t\"%s\": struct{}{},\n" $(echo $pkg | tr -d \\r) >> $dest
    done

    # whats removed in new release not really required.
    # printf "(removed) Go v1.%s vs v1.%s\n" $i "$(($i-1))"
    # comm -13  <(sort go@1.$i) <(sort go@1.$(($i-1)))
done

echo "}" >> $dest


gofmt -w $dest
