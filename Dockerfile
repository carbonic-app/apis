FROM golang:1.13 as builder

WORKDIR /go/src/app
ADD . .

RUN go mod download
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -v -o /main ./pkg/cmd/server

FROM scratch
COPY --from=builder /main ./
CMD ["./main", "-grpc-port", "80", "-db-file", "./user.db"]
EXPOSE 80
