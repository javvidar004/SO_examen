FROM golang:1.23.8

WORKDIR /app

# Install dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download
# Copy the rest of the application code
COPY main.go .


CMD [ "go", "run", "main.go" ]