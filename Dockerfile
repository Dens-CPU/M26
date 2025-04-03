FROM golang
RUN mkdir -p /go/src/M26
WORKDIR /go/src/M26
ADD go.mod .
ADD main.go .
RUN go install .

FROM alpine:latest
LABEL version="v1.0"
LABEL maintainer="Dens"
WORKDIR /root/
COPY --from=0 /go/bin/M26 . 
ENTRYPOINT ./M26
