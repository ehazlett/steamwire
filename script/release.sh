#!/bin/bash
OS="linux windows"
ARCH="386 amd64"
TMPDIR="build"
VERSION=$(cat version/version.go | grep "version     = " | tr -d '[:space:]' | sed 's/version=//g' | tr -d '"')

echo "Building release for version ${VERSION}"

for x in $OS; do
    for a in $ARCH; do
        echo " -> building ${x} ${a}"
        mkdir -p ${TMPDIR}/${x}/${a}
        make GOOS=$x GOARCH=$a build > /dev/null 2>&1
        if [ $? -ne 0 ]; then
            echo " ERR: error building for ${x} ${a}"
            continue
        fi
        mv ./cmd/steamwire/steamwire* ${TMPDIR}/${x}/${a}/
        zip -D -r ${TMPDIR}/steamwire-${VERSION}-${x}-${a}.zip ${TMPDIR}/${x}/${a}/ > /dev/null 2>&1
        if [ $? -ne 0 ]; then
            echo " ERR: packaging build for ${x} ${a}"
            continue
        fi
    done

    # cleanup
    rm -rf ${TMPDIR}/${x}
done
