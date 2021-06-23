.PHONY: clean
DIST_DIR="./cache/build"

clean:
	@rm -rf build
	# @rm -rf ${DIST_DIR}

build: clean
	@mkdir -p ${DIST_DIR}
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build \
		-ldflags "-X main.Version=$(VERSION)" \
		-o ${DIST_DIR}/cache cache/main.go