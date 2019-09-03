FROM golang:alpine as builder
WORKDIR /app
COPY go.mod .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go test -v ./test && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a \
    -installsuffix cgo -o ./bin/server ./cmd/http-server/

FROM scratch
WORKDIR /app
COPY --from=builder /app/bin/server .
EXPOSE 8081
ENTRYPOINT ["./server"]