FROM golang:alpine as builder

WORKDIR /src

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o pastebin .

FROM alpine:3.14

WORKDIR /opt/app

COPY --from=builder /src/pastebin /opt/app/

ENV GIN_MODE=release

EXPOSE 8080

CMD ["/opt/app/pastebin"]
