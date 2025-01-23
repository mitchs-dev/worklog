#!/bin/bash

symVer=$1
commit=$(git rev-parse HEAD)
buildDate=$(date -u +'%Y-%m-%dT%H:%M:%SZ')
versionFile="internal/version/version.json"
binDir="bin"

if [ -z "$symVer" ]; then
  echo "Semver is required as the first argument"
  exit 1
fi

# If the version file doesn't exist, create it
if [ ! -f $versionFile ]; then
  echo "Creating version file ($versionFile) with symVer $symVer, commit $commit and build date $buildDate"
  mkdir -p internal/version
  echo "{\"semantic\": \"$symVer\", \"build\": {\"commit\": \"$commit\", \"date\": \"$buildDate\"}}" > $versionFile
fi

echo "Updating version file ($versionFile) with symVer $symVer, commit $commit and build date $buildDate"
jq ".semantic = \"$symVer\"" $versionFile > tmp.json && mv tmp.json $versionFile
jq ".build.commit = \"$commit\"" $versionFile > tmp.json && mv tmp.json $versionFile
jq ".build.date = \"$buildDate\"" $versionFile > tmp.json && mv tmp.json $versionFile


echo "Building binaries for version $symVer with hash $commit and build date $buildDate"
rm -rf $binDir
mkdir -p $binDir

for os in "linux" "darwin" "windows"; do
  for arch in "amd64" "arm64"; do
    echo "Building $os $arch"
    if [ "$os" == "windows" ]; then
      # Skip windows arm64 (not supported)
      if [ "$arch" == "arm64" ]; then
        continue
      fi
      binDir="bin/windows"
      binName="worklog-v$symVer.exe"
    else
      binDir="bin/$os"
      binName="worklog-$os-$arch-v$symVer"
    fi
    GOOS=$os GOARCH=$arch go build -ldflags="-s -w -buildid=" -o $binDir/$binName ./cmd/main.go
    shasum -a 256 $binDir/$binName > $binDir/$binName.sha256
    chmod +x $binDir/$binName
  done
done
