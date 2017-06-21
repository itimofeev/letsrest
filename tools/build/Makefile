dep:
	godep restore ./...

build: dep
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s" -a -installsuffix cgo -o /letsrest github.com/itimofeev/letsrest/main/letsrest
