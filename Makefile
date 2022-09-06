PACKAGE=github.com/nikhil1raghav/kindle-send
VERSION=$(shell cat VERSION)
BUILD_DATE=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
CURRENT_DIR=$(shell pwd)
DIST_DIR=${CURRENT_DIR}/dist

override LDFLAGS += \
-s -w \
-X ${PACKAGE}/util.version=${VERSION}\
-X ${PACKAGE}/util.buildDate=${BUILD_DATE}

linux:
	 CGO=0 GOOS=linux GOARCH=amd64 go build -ldflags '${LDFLAGS}' -o ${DIST_DIR}/kindle-send-linux-64bit-${VERSION} ./main.go
	 CGO=0 GOOS=linux GOARCH=386 go build -ldflags '${LDFLAGS}' -o ${DIST_DIR}/kindle-send-linux-32bit-${VERSION} ./main.go
	 CGO=0 GOOS=linux GOARCH=arm go build -ldflags '${LDFLAGS}' -o ${DIST_DIR}/kindle-send-linux-arm-${VERSION} ./main.go
	 CGO=0 GOOS=linux GOARCH=arm64 go build -ldflags '${LDFLAGS}' -o ${DIST_DIR}/kindle-send-linux-arm64-${VERSION} ./main.go
	 upx -7 ${DIST_DIR}/kindle-send-linux-64bit-${VERSION}
	 upx -7 ${DIST_DIR}/kindle-send-linux-arm64-${VERSION}
	 upx -7 ${DIST_DIR}/kindle-send-linux-32bit-${VERSION}
	 upx -7 ${DIST_DIR}/kindle-send-linux-arm-${VERSION}


#not packing windows binary, defender flags upx packed binary as trojan :( 
windows:
	CGO=0 GOOS=windows GOARCH=amd64 go build -ldflags '${LDFLAGS}' -o ${DIST_DIR}/kindle-send-windows-64bit-${VERSION}.exe ./main.go
	CGO=0 GOOS=windows GOARCH=386 go build -ldflags '${LDFLAGS}' -o ${DIST_DIR}/kindle-send-windows-32bit-${VERSION}.exe ./main.go
	CGO=0 GOOS=windows GOARCH=arm go build -ldflags '${LDFLAGS}' -o ${DIST_DIR}/kindle-send-windows-arm-${VERSION}.exe ./main.go

darwin:
	CGO=0 GOOS=darwin GOARCH=amd64 go build -ldflags '${LDFLAGS}' -o ${DIST_DIR}/kindle-send-darwin-64bit-${VERSION} ./main.go
	CGO=0 GOOS=darwin GOARCH=arm64 go build -ldflags '${LDFLAGS}' -o ${DIST_DIR}/kindle-send-darwin-arm64-${VERSION} ./main.go
