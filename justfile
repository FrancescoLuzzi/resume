default:
    @just --list

build:
    @go build

generate:
    #!/usr/bin/env bash
    rm -rf build
    mkdir build
    for res in $(ls locales);do
        lang=$(echo $res | cut -f1 -d.)
        ./resume -cv locales/$lang.yaml -author author.yaml --out $lang.typ
        typst compile $lang.typ --font-path ./fonts --ignore-system-fonts build/$lang.pdf
        rm $lang.typ
    done
