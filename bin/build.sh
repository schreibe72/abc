rm *.zip
rm -rf darwin
rm -rf linux
rm -rf windows

VERSION=$(git name-rev --tags --name-only $(git rev-parse HEAD))
GITHASH=$(git log -n1 --pretty='%h')

mkdir -p darwin/amd64
mkdir -p linux/amd64
mkdir -p linux/i386
mkdir -p windows/amd64
mkdir -p windows/i386

echo VERSION: $VERSION
echo GITHASH: $GITHASH
cd darwin/amd64
env GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.version=$VERSION -X main.githash=$GITHASH" ../../../
zip ../../darwin_amd64.zip abc
cd ../../linux/amd64
env GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=$VERSION -X main.githash=$GITHASH" ../../../
zip ../../linux_amd64.zip abc
cd ../i386
env GOOS=linux GOARCH=386 go build -ldflags "-X main.version=$VERSION -X main.githash=$GITHASH" ../../../
zip ../../linux_i386.zip abc
cd ../../windows/amd64
env GOOS=windows GOARCH=amd64 go build -ldflags "-X main.version=$VERSION -X main.githash=$GITHASH" ../../../
zip ../../windows_amd64.zip abc.exe
cd ../i386
env GOOS=windows GOARCH=386 go build -ldflags "-X main.version=$VERSION -X main.githash=$GITHASH" ../../../
zip ../../windows_i386.zip abc.exe
cd ../../

cp darwin/amd64/abc $GOPATH/bin
abc save bash -f ~/abc.completion.sh
