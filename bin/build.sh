rm *.zip
rm -rf darwin
rm -rf linux
rm -rf windows

mkdir -p darwin/amd64
mkdir -p linux/amd64
mkdir -p linux/i386
mkdir -p windows/amd64
mkdir -p windows/i386

cd darwin/amd64
env GOOS=darwin GOARCH=amd64 go build ../../../
zip ../../darwin_amd64.zip abc
cd ../../linux/amd64
env GOOS=linux GOARCH=amd64 go build ../../../
zip ../../linux_amd64.zip abc
cd ../i386
env GOOS=linux GOARCH=386 go build ../../../
zip ../../linux_i386.zip abc
cd ../../windows/amd64
env GOOS=windows GOARCH=amd64 go build ../../../
zip ../../windows_amd64.zip abc.exe
cd ../i386
env GOOS=windows GOARCH=386 go build ../../../
zip ../../windows_i386.zip abc.exe
cd ../../
