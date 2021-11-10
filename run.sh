#!/bin/bash
set -Eeux

SCRIPT_DIR="$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
mkdir -p /tmp/juscaba/pdfs
mkdir -p /tmp/juscaba/web/public/data
for ((i=1; i<=$#; i++))
do
    exp=${!i}
    exp_filename=${!i/\//-}

    ./builder "-json=/tmp/juscaba/web/public/data/${exp_filename}.json" -pdfs=/tmp/juscaba/pdfs "-expediente=${exp}" -images=false "-blacklist=${BLACKLIST_REGEX:-}"

    pushd ts
    yarn run ts-node create-index.ts /tmp/juscaba/web/public/data/${exp_filename}.json /tmp/juscaba/web/public/data/${exp_filename}-index.json
    popd
done

pushd web
yarn
yarn build
popd

mkdir -p /tmp/juscaba/web
rm -rf /tmp/juscaba/web/build
chmod -R 777 build
mv build /tmp/juscaba/web/
