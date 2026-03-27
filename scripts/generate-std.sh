#!/usr/bin/env bash

mkdir -p .tmp

tip=26

# create std pkg list

docker run --rm -it golang:1.${tip}-alpine go list std | grep -v internal | grep -v vendor > .tmp/go@1.${tip};

dest="analyzer/std.go"

echo "// Code generated using std.sh; DO NOT EDIT." > $dest
echo "" >> $dest

echo "// We will ignore that fact that some of packages" >> $dest
echo "// were removed from stdlib." >> $dest
echo "" >> $dest
echo "package analyzer" >> $dest
echo "" >> $dest
echo "var std = map[string]struct{}{" >> $dest

for pkg in $(sort .tmp/go@1.${tip}); do
    printf "\t\"%s\":  {},\n" "$(echo "$pkg" | tr -d \\r)" >> $dest
done

echo "}" >> $dest


gofmt -w $dest
