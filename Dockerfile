FROM golang:1.26-alpine AS builder
ENV GOTOOLCHAIN=local
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o bot .

FROM gcr.io/distroless/static-debian12
COPY --from=builder /app/bot /bot
ENTRYPOINT ["/bot"]
