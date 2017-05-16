export GOOS=linux
export GOARCH=mips32
export GOROOT=/opt/mipsgo
export PATH=/opt/mipsgo/bin:$PATH
echo "Begin to build"
go build
echo "All done~"
