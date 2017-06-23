FROM golang:1.8.3

RUN go get github.com/tools/godep && mkdir /_goTestOutput

WORKDIR /
ADD tools/build/Makefile /
COPY ./ /usr/local/go/src/github.com/itimofeev/letsrest

CMD ["make", "build"]
