FROM golang:latest

ENV GOPATH=/

RUN apt-get update
RUN apt-get -y install postgresql-client

COPY services/ services

RUN chmod +x services/tweet/wait-for-postgres.sh

COPY go.work .

COPY services/tweet/go.mod services/tweet

RUN cd services/tweet && go mod tidy && go mod download && go clean --modcache

RUN cd services/tweet && go build -o twitter-clone-tweet ./cmd/main.go

EXPOSE 7006

CMD [ "./main" ]