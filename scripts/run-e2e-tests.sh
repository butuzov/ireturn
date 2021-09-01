#!/usr/bin/env bash

pushd .sandbox;

    for repo in $(cat repos.txt)
    do
        name=$(basename $repo);
        if [[ ! -d "repo-$(basename $repo)" ]]; then
            git clone -q $repo "repo-$(basename $repo)";
        fi
    done


    for d in repo-*/; do
        pushd $d;
            printf "\n ==== [%s] ========= \n" $d;
            time ireturn ./...
        popd;
    done

popd;
