export GOOS=linux
export GOARCH=arm
export GOARM=5
# export GOROOT=/opt/mipsgo
#export PATH=/opt/mipsgo/bin:$PATH
echo "Begin to build"
go build
echo "All done~"
