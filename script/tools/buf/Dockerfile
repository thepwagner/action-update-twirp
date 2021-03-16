FROM alpine:3.13.2 AS buf

ARG BUF_VERSION=v0.40.0
ARG BUF_CHECKSUM=ce9b0bea308970d6ac023ed9a6fb45e77647575ca4ff5b121df898243d60a519

ARG BUFF_URL=https://github.com/bufbuild/buf/releases/download/${BUF_VERSION}/buf-Linux-x86_64
RUN apk --no-cache add --virtual .build curl \
  && curl -o /usr/local/bin/buf -L "$BUFF_URL" \
  && echo "$BUF_CHECKSUM  /usr/local/bin/buf" | sha256sum -c \
  && chmod +x /usr/local/bin/buf \
  && apk del .build


FROM golang:1.16.2 AS protoc-plugins-go
RUN mkdir /app
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
RUN go build -o /protoc-gen-go github.com/golang/protobuf/protoc-gen-go
RUN go build -o /protoc-gen-twirp github.com/twitchtv/twirp/protoc-gen-twirp


FROM buf AS protoc
COPY --from=protoc-plugins-go /protoc-gen-go /usr/local/bin/protoc-gen-go
COPY --from=protoc-plugins-go /protoc-gen-twirp /usr/local/bin/protoc-gen-twirp