bash 
FROM golang:1.23-rc-alpine
WORKDIR /app
RUN apk add --no-cache postgresql-client
COPY go.mod ./
COPY go.sum ./


RUN go mod download
COPY . .
RUN go build -o main .
EXPOSE 8080

CMD ["./main"]
