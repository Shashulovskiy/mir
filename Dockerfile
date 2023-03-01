FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
CMD go run ./samples/brb-channel
#CMD [". /node"]