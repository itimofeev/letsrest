FROM golang:1.8.3

RUN go get github.com/tools/godep

WORKDIR /
ADD tools/build/Makefile /
COPY ./ /go/src/github.com/itimofeev/letsrest/

CMD ["make", "build"]
