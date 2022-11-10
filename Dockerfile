FROM golang:latest

ENV GOPATH=/

COPY ./ ./

RUN apt-get update
RUN apt-get -y install postgresql-client

RUN chmod +x wait-for-postgres.sh

COPY go.mod go.sum ./
RUN go mod download && go mod verify

RUN go build -o twitter-clone ./cmd/main.go

EXPOSE 7000

CMD [ "./main" ]