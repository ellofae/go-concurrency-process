FROM golang:1.21-alpine as builder

RUN apk update --no-cache && apk upgrade --no-cache

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -o /go/app ./cmd/app

FROM scratch as build
COPY --from=builder /bin/app /bin/app

CMD ["/bin/app"]