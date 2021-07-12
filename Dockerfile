######## Start a builder stage #######
FROM golang:1.16-alpine as builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/main cmd/auth-service/main.go

######## Start a new stage from scratch #######
FROM alpine:latest

RUN apk --no-cache add tzdata
WORKDIR /opt/
COPY --from=builder /app/bin/main .

CMD ["./main"]