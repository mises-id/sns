FROM golang:1.16-alpine3.14 AS builder
RUN apk update && apk add --no-cache git
ENV GOPRIVATE=github.com/mises-id

WORKDIR /opt/mises
COPY go.mod go.sum ./
RUN go mod download

COPY . .
# Mark the build as statically linked.
RUN CGO_ENABLED=0 go build -o /bin/mises cmd/main.go

# FROM scratch AS final
FROM alpine:3.14 AS final

WORKDIR /opt/mises
COPY --from=builder /bin/mises /bin/

ENTRYPOINT ["/bin/mises"]
EXPOSE 8080
