FROM golang:alpine as builder
RUN apk add --no-cache git
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o ./bin/server ./cmd/http-server/

FROM scratch
WORKDIR /app
COPY --from=builder /app/bin/server .
EXPOSE 8081
ENTRYPOINT ["./server"]