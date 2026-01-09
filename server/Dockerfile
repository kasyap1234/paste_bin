FROM golang:1.25.5-alpine as builder

RUN apk add --no-cache git ca-certificates tzdata
RUN adduser -D -g '' appuser
WORKDIR /build
COPY go.mod  go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /app/server ./cmd/pastebin-api
FROM scratch
COPY --from=builder /app/server /server
COPY --from=builder /etc/ssl/certs/ca-certificates.crt  /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
ENTRYPOINT [ "/server" ]
