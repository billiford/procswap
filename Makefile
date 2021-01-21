build:
		go build -o bin/procswap cmd/procswap/procswap.go

run:
		go run cmd/procswap/procswap.go

compile:
		echo "Compiling for every OS and Platform"
		GOOS=linux GOARCH=arm go build -o bin/procswap-linux-arm cmd/procswap/procswap.go
		GOOS=linux GOARCH=arm64 go build -o bin/procswap-linux-arm64 cmd/procswap/procswap.go
		GOOS=freebsd GOARCH=386 go build -o bin/procswap-freebsd-386 cmd/procswap/procswap.go
		GOOS=windows go build -o bin/procswap.exe cmd/procswap/procswap.go
