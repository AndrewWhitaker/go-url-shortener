FROM golang:1.18-alpine AS build
WORKDIR /app
RUN apk add git

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o /url-shortener


FROM alpine:3.15.4
RUN apk update && apk add bash
WORKDIR /

COPY --from=build /url-shortener /url-shortener

EXPOSE 8080

CMD ["./url-shortener"]

