# This Dockerfile is used for testing and local development. The 'release' image
# used by the CI system is within the '.build' directory,
FROM golang:1.14.0
ENV GOARCH=amd64 GOOS=linux
WORKDIR /go/src/github.com/fernandrone/linelint/
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /bin/linelint

FROM scratch
COPY --from=0 /bin/linelint /linelint
COPY LICENSE README.md ./
WORKDIR /data
ENTRYPOINT ["/linelint", "."]
