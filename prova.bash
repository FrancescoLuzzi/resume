for res in $(ls ./locales); do
    echo round 1
    lang=$(echo $lang | cut -f1 -d.)
    echo $lang
    ./resume -cv locales/$lang.yaml -author author.yaml --out $lang.typ
    echo "resume done"
    typst compile $lang.typ --font-path ./fonts --ignore-system-fonts build/$lang.pdf
    echo "compiled"
    rm $lang.typ
done
