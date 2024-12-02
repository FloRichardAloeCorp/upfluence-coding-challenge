#### BUILD STAGE
FROM golang:1.22.4-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o upfluence-coding-challenge

#### IMAGE DEFINITION
FROM scratch

COPY --chown=1000:1000 --from=builder /app/upfluence-coding-challenge /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

USER 1000:1000

CMD ["./upfluence-coding-challenge"]