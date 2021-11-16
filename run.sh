#!/bin/bash
set -Eeux

mkdir -p /tmp/juscaba/pdfs
mkdir -p public/data
for exp in $EXPEDIENTES
do
    exp_filename=${!i/\//-}

    ./builder "-json=public/data/${exp_filename}.json" -pdfs=/tmp/juscaba/pdfs "-expediente=${exp}" -images=${READ_IMAGES:-true} "-blacklist=${BLACKLIST_REGEX:-}"

    pushd ts
    yarn run ts-node create-index.ts ../public/data/${exp_filename}.json ../public/data/${exp_filename}-index.json
    popd
done


sed -i "s~%HOMEPAGE%~${URL}~g" web/package.json
yarn build

rm -rf /tmp/juscaba/build
chmod -R 777 build
mv build /tmp/juscaba/
