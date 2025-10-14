FROM golang:alpine AS build_container

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o server

FROM alpine:3.18

WORKDIR /app

COPY --from=build_container /app/server .

EXPOSE 8080

CMD ["./server"]


