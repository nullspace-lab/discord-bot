FROM golang:1.26.3-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -v -o bot .

FROM gcr.io/distroless/static-debian12
COPY --from=builder /app/bot /bot
ENTRYPOINT ["/bot"]
