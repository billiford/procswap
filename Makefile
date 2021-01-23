build:
		go build -o bin/procswap cmd/procswap/procswap.go

run:
		go run cmd/procswap/procswap.go

REVISION := $(shell git rev-parse HEAD)

compile:
		echo "Compiling for every OS and Platform"
		GOOS=linux GOARCH=arm go build -o bin/procswap-linux-arm -ldflags "-X main.version=$(VERSION) -X main.revision=$(REVISION)" cmd/procswap/procswap.go
		GOOS=linux GOARCH=arm64 go build -o bin/procswap-linux-arm64 -ldflags "-X main.version=$(VERSION) -X main.revision=$(REVISION)" cmd/procswap/procswap.go
		GOOS=linux GOARCH=amd64 go build -o bin/procswap-linux-amd64 -ldflags "-X main.version=$(VERSION) -X main.revision=$(REVISION)" cmd/procswap/procswap.go
		GOOS=freebsd GOARCH=386 go build -o bin/procswap-freebsd-386 -ldflags "-X main.version=$(VERSION) -X main.revision=$(REVISION)" cmd/procswap/procswap.go
		GOOS=windows go build -o bin/procswap.exe -ldflags "-X main.version=$(VERSION) -X main.revision=$(REVISION)" cmd/procswap/procswap.go
