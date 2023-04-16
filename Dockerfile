FROM golang:1.18-alpine

RUN apk add build-base

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

CMD GOOS=linux CGO_CFLAGS="-O2" CGO_CXXFLAGS="-std=c++11" go run ./samples/brb-channel tests none
#CMD [". /node"]