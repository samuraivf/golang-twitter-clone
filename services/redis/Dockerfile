FROM golang:latest

ENV GOPATH=/

RUN apt-get update
RUN apt-get -y install redis

COPY services/ services

COPY go.work .

COPY services/redis/go.mod services/redis
COPY services/redis/go.sum services/redis

RUN cd services/redis && go mod tidy && go mod download && go clean --modcache

RUN cd services/redis && go build -o twitter-clone-redis ./cmd/main.go

EXPOSE 7002

CMD [ "services/redis/twitter-clone-redis" ]