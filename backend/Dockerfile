FROM golang:1.22.2

WORKDIR /

COPY . .

# ENV GO111MODULE=off

RUN go mod tidy
RUN go build -o bin/main ./cmd

EXPOSE 8080

CMD ["/bin/main"]