# This Dockerfile is used for testing and local development. The 'release' image
# used by the CI system is within the '.build' directory,
FROM golang:1.14.0
WORKDIR /go/src/github.com/fernandrone/linelint/
COPY go.mod go.sum ./
RUN go get ./...
COPY . .
RUN GOARCH=amd64 GOOS=linux go build -o /bin/linelint

FROM scratch
COPY /bin/linelint /linelint
COPY LICENSE README.md ./
ENTRYPOINT ["/linelint"]
