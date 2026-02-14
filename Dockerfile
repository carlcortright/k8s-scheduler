FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /scheduler ./cmd

FROM alpine:3.20
RUN apk --no-cache add ca-certificates
COPY --from=builder /scheduler /scheduler
USER nobody
ENTRYPOINT ["/scheduler"]
