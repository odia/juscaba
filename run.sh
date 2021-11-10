#!/bin/bash
set -Eeux

SCRIPT_DIR="$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
mkdir -p /tmp/juscaba/pdfs
mkdir -p public/data
for ((i=1; i<=$#; i++))
do
    exp=${!i}
    exp_filename=${!i/\//-}

    ./builder "-json=public/data/${exp_filename}.json" -pdfs=/tmp/juscaba/pdfs "-expediente=${exp}" -images=false "-blacklist=${BLACKLIST_REGEX:-}"

    pushd ts
    yarn run ts-node create-index.ts public/data/${exp_filename}.json public/data/${exp_filename}-index.json
    popd
done

yarn
yarn build

mkdir -p /tmp/juscaba/web
rm -rf /tmp/juscaba/web/build
chmod -R 777 build
mv build /tmp/juscaba/web/
