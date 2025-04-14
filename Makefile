GOFLAGS = -trimpath -ldflags='-s -w -H=windowsgui'
export GOOS = windows

build:
	go build ${GOFLAGS} -o bin/fastbrowser.exe