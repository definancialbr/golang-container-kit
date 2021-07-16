FROM golang:1.16-buster

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY *.go .

RUN cd cmd/noop && go build -o /noop

CMD ["/noop"]