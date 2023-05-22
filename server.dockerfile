# base image
FROM golang:1.20-alpine as server-builder
RUN mkdir /app
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 go build -o server ./cmd/server
RUN chmod +x /app/server

# small image with just executable
FROM alpine:latest
RUN mkdir /app
COPY --from=server-builder /app/server /app
CMD ["/app/server"]