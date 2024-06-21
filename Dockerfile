FROM golang:1.19-bullseye AS builder

WORKDIR /app

COPY . .

RUN go mod download

# I had some trouble using the -o flag for some reason.
# We can't set CGO_ENABLED=0 at least because of SQLite.
RUN go build github.com/appditto/pippin_nano_wallet/apps/cli

FROM debian:stable-slim

WORKDIR /root/

COPY --from=builder /app/cli .

CMD ["./cli", "--start-server"]
