FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /app ./cmd/app
FROM alpine:3.23
RUN adduser -D -u 10001 app
USER app
COPY --from=build /app /app
EXPOSE 8080
ENTRYPOINT ["/app"]
