FROM golang:1.18-alpine

RUN apk add build-base

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

CMD go env -w CGO_CXXFLAGS="-std=c++11 -O2"
CMD go run ./samples/brb-channel tests corrupt
#CMD [". /node"]