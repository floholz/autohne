# syntax=docker/dockerfile:1

FROM golang:1.21

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY ../go.mod ../go.sum ./
RUN go mod download

# Copy the source code.
COPY ../app.go ./
COPY ../src/* ./src/
COPY ../assets/* ./asstes/

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /autohne

# Default TCP port
EXPOSE 7887

# Run
CMD ["/autohne"]