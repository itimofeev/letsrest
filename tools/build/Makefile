dep:
	cd /usr/local/go/src/github.com/itimofeev/letsrest ; GOPATH=/go godep restore -v ./...

test: dep
	cd /usr/local/go/src/github.com/itimofeev/letsrest ; go test -v ./... | tee /_goTestOutput/goTestOutput0.log

build: dep
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s" -a -installsuffix cgo -o /letsrest github.com/itimofeev/letsrest/main/letsrest
