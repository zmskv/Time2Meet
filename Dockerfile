FROM golang:1.25.1-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/time2meet ./cmd/api

FROM alpine:3.20
RUN adduser -D -H -u 10001 appuser
USER appuser

COPY --from=build /bin/time2meet /usr/local/bin/time2meet

EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/time2meet"]

