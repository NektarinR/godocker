FROM golang:alpine as builder
RUN apk add --no-cache git
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go test ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/server ./cmd/http-server/

FROM scratch
WORKDIR /app
COPY --from=builder /app/bin/server .
ENTRYPOINT ["./server"]