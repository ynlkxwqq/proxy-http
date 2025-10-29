FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o proxy-server main.go

EXPOSE 3333

CMD ["./proxy-server"]

