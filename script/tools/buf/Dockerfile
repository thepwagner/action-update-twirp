FROM alpine:3.14.0 AS buf

ARG BUF_VERSION=v0.43.2
ARG BUF_CHECKSUM=53933427cb5518830a3fcf11d3bd9b6e3fb36eae9086731484b8ab50630c9d09

ARG BUFF_URL=https://github.com/bufbuild/buf/releases/download/${BUF_VERSION}/buf-Linux-x86_64
RUN apk --no-cache add --virtual .build curl \
  && curl -o /usr/local/bin/buf -L "$BUFF_URL" \
  && echo "$BUF_CHECKSUM  /usr/local/bin/buf" | sha256sum -c \
  && chmod +x /usr/local/bin/buf \
  && apk del .build


FROM golang:1.16.5 AS protoc-plugins-go
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