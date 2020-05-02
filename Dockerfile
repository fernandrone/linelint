FROM golang:1.14
WORKDIR /go/src/github.com/fernandrone/linelint/
ADD go.mod go.sum ./
RUN go get ./...
ADD . .
RUN GOOS=linux CGO_ENABLED=0 go build -o /bin/linelint

FROM scratch
COPY --from=0 /bin/linelint /linelint
ADD LICENSE README.md ./
ENTRYPOINT ["/linelint"]
