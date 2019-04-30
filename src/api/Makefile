NAME    := articles
SRCS    := $(shell find . -type f -name '*.go')
LDFLAGS := -ldflags="-extldflags \"-static\""
GBFLAGS := -a -installsuffix netgo

.PHONY: fmt local product clean bindata-local bindata

fmt:
	go fmt

bindata-local:
	go-bindata -debug=true -tags=local -o=bindata_local.go ./template/... ./static/...

bindata:
	go-bindata -nomemcopy -tags=!local -o=bindata.go ./template/... ./static/...

local: fmt bindata-local $(SRCS)
	go build ${LDFLAGS} ${GBFLAGS} -tags="local netgo" -o=bin/${NAME}-local

product: fmt bindata $(SRCS)
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} ${GBFLAGS} -tags="product netgo" -o=bin/${NAME}

clean:
	rm -rf bin/*

clean-bindata:
	rm -f bindata.go bindata_local.go

build-test: clean local
