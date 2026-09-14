FROM golang:1.25-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /app/bin/backend ./cmd/backend

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/bin/backend ./backend
EXPOSE 8080
ENTRYPOINT ["./backend"]
