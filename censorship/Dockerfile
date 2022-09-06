FROM golang:1.18-alpine

ENV GOPATH=/

COPY ./ ./


RUN go mod download
RUN go build -o censorship ./cmd/main.go


CMD ["./censorship"]