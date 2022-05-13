# Backend (Go) builder
FROM golang:1.18-alpine AS backend
WORKDIR /app
RUN apk add git

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o /url-shortener

# Frontend (ReactJS) builder
FROM node:16 as frontend
WORKDIR /app
COPY assets/ .

WORKDIR /app/assets
RUN npm install && npm run build
RUN ls -lah

# Final image
FROM alpine:3.15.4
RUN apk update && apk add bash
WORKDIR /

COPY --from=backend /url-shortener /url-shortener
COPY --from=frontend /app/build/ ./assets/build

EXPOSE 8080

CMD ["./url-shortener"]

