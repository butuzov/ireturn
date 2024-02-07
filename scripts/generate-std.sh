#!/usr/bin/env bash

mkdir -p .tmp

tip=22

# create std pkg list
for i in $(seq 1 $tip); do
    if [[ ! -a .tmp/go@1.$i ]]; then
        docker run --rm -it golang:1.$i go list std \
            | grep -v internal | grep -v vendor > .tmp/go@1.$i;
        docker rmi golang:1.$i
    fi
done

dest="analyzer/std.go"

echo "// Code generated using std.sh; DO NOT EDIT." > $dest
echo "" >> $dest

echo "// We will ignore that fact that some of packages" >> $dest
echo "// were removed from stdlib." >> $dest
echo "" >> $dest
echo "package analyzer" >> $dest
echo "" >> $dest
echo "var std = map[string]struct{}{" >> $dest

for i in $(seq 2 $tip); do

    printf "\t// added in Go v1.%s in compare to v1.%s (docker image)\n" $i "$(($i-1))" >> $dest

    for pkg in $(comm -13 <(sort .tmp/go@1.$(($i-1))) <(sort .tmp/go@1.$i)); do
        printf "\t\"%s\":  {},\n" $(echo $pkg | tr -d \\r) >> $dest
    done

done

echo "}" >> $dest


gofmt -w $dest
