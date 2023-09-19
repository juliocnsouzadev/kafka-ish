# base image
FROM golang:latest as server-builder
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 go build -o kafka-ish-server ./cmd/server
RUN chmod +x /app/kafka-ish-server

# small image with just executable
FROM debian:latest
RUN echo "building server executable"
RUN mkdir /app
COPY --from=server-builder /app/kafka-ish-server /app
CMD ["/app/kafka-ish-server"]