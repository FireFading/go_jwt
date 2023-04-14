FROM golang:latest

WORKDIR /app

# Copy the Go mod and sum files into the working directory
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code from the current directory to the working directory inside the container
COPY . .

# Build the Go app
RUN go build -o app .

EXPOSE 3000

