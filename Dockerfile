FROM golang:1.21-alpine as builder

RUN --mount=type=cache,target=/go/pkg/mod/ \
     apk update && apk upgrade

WORKDIR /src

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x
COPY . .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -o /cmd/app ./cmd/app

FROM busybox as intermediate

RUN mkdir /config
COPY --from=builder /src/config/config.yaml /config/config.yaml

FROM scratch as build

COPY --from=builder /cmd/app /bin/app
COPY --from=intermediate /config/config.yaml /config/config.yaml

ENTRYPOINT ["/bin/app"]