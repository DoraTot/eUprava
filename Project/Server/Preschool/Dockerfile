##FROM golang:alpine AS build_container
##WORKDIR /app
##COPY go.mod .
###COPY go.sum .
##RUN go mod download
##COPY . .
##RUN go build -o server
##
##FROM alpine:3.18
##
##EXPOSE 8080
##
###ENTRYPOINT ["server"]
##CMD ["./main"]
#
#FROM golang:1.20
#
#WORKDIR /app
#
##COPY go.mod go.sum ./
#COPY go.mod ./
#
#RUN go mod download
#
#COPY . .
#
#RUN go build -o main .
#
#EXPOSE 8080
#
#CMD ["./main"]
#

# Build Stage
FROM golang:alpine AS build_container

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o server

# Runtime Stage
FROM alpine:3.18

WORKDIR /app

COPY --from=build_container /app/server .

EXPOSE 8080

CMD ["./server"]


#FROM golang:alpine AS build_container
#WORKDIR /app
#COPY go.mod .
##COPY go.sum .
#RUN go mod download
#COPY . .
#RUN go build -o main
##RUN go build -o main .
#
#FROM alpine:3.18
#
#
#EXPOSE 8080
#
##ENTRYPOINT ["server"]
#CMD ["./main"]


