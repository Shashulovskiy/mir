FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

CMD GOOS=linux GOARCH=amd64 go run ./samples/brb-channel tests none
#CMD [". /node"]