FROM golang:1.18-alpine AS builder

WORKDIR /async-challenge

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
# Expecting to copy go.mod and if present go.sum.
COPY go.* ./
RUN go mod download

# Copy local code to the container image.
COPY .. ./

# Build the binary.
RUN go build -ldflags="-s -w" -v -o queueprocessor ./cmd/

FROM alpine:latest AS runtime

COPY --from=builder /mallbots/docker/wait-for .
RUN chmod +x /wait-for

# Copy the binary to the production image from the builder stage.
COPY --from=builder /async-challenge/queueprocessor /async-challenge/queueprocessor

CMD ["/async-challenge/queueprocessor"]
