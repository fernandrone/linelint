# This Dockerfile is used by the CI system.
FROM scratch
COPY bin/linelint /linelint
COPY LICENSE README.md ./
WORKDIR /data
ENTRYPOINT ["/linelint", "."]
