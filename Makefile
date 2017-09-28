BINARY      = terraform-provider-vbox
GOARCH      = amd64
VERSION     = 0.0.1
CURRENT_DIR = $(shell pwd)
LDFLAGS     = -ldflags "-X main.VERSION=${VERSION}"

# Build the project
all: clean test linux

linux:
	cd ${CURRENT_DIR}; \
	GOOS=linux GOARCH=${GOARCH} go build ${LDFLAGS} -o ${GOBIN}/${BINARY} . ; \
  mkdir -p ~/.terraform.d/plugins/linux_amd64 ; \
  cp ${GOBIN}/${BINARY} ~/.terraform.d/plugins/linux_amd64/${BINARY}

clean:
	rm -f ${BINARY} ${GOBIN}/${BINARY}

.PHONY: linux test clean
