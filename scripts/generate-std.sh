#!/usr/bin/env bash

mkdir -p .tmp

tip=26

# create std pkg list

docker run --rm -it golang:1.${tip}-alpine go list std | grep -v internal | grep -v vendor > .tmp/go@1.${tip};

dest="analyzer/std.go"

cat > $dest <<EOF
// Code generated using scripts/generate-std.sh; DO NOT EDIT.

// We will ignore that fact that some of packages
// were removed from stdlib.

package analyzer

var std = map[string]struct{}{
EOF

for pkg in $(sort .tmp/go@1.${tip}); do
    printf "\t\"%s\":  {},\n" "$(echo "$pkg" | tr -d \\r)" >> $dest
done

echo "}" >> $dest


gofmt -w $dest
