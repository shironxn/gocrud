FROM golang:1.22.2

# WORKDIR /

# COPY . .

# # ENV GO111MODULE=off

# RUN go mod tidy
# RUN go build -o bin/main ./cmd

# EXPOSE 8080

# CMD ["/bin/main"]

# Start from the latest golang base image

# Add Maintainer Info
LABEL maintainer="qshiro69@gmail.com"

# Set the Current Working Directory inside the container
WORKDIR /

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

RUN go mod tidy
# Disable Go Modules
ENV GO111MODULE=off

# Build the Go app
RUN go build -o main ./cmd

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]