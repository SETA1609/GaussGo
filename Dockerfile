FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /gaussgo ./cmd/gaussgo

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /gaussgo .
COPY --from=builder /app/data ./data
COPY --from=builder /app/pdfs ./pdfs

CMD ["./gaussgo"]