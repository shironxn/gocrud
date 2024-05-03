FROM golang:1.22.2

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o bin/main ./cmd

EXPOSE ${APP_PORT}

CMD ["/app/bin/main"]