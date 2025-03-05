FROM golang:1.23.6-alpine

WORKDIR /goSMA

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o goSMA

ENTRYPOINT ["./goSMA"]
