FROM golang
RUN mkdir -p /go/src/M26
WORKDIR /go/src/M26
ADD main.go .
ADD go.mod .
RUN go install .

FROM alpine:latest
LABEL version="v1.0"
LABEL maintainer="Dens"
WORKDIR /root/
COPY --from=0 /go/bin/pipline .
ENTRYPOINT ./pipline

