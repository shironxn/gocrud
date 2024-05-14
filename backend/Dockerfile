FROM golang:1.22.2

WORKDIR /app

COPY . .

ENV GO111MODULE=off

RUN go mod download
RUN go build -o bin/main ./cmd

EXPOSE 8080

CMD ["/app/bin/main"]