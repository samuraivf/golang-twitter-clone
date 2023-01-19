FROM golang:1.19

ENV GOPATH=/

RUN apt-get update
RUN apt-get -y install postgresql-client

COPY services/ services

RUN chmod +x services/user/wait-for-postgres.sh

COPY go.work .

COPY services/user/go.mod services/user
COPY services/user/go.sum services/user

RUN cd services/user && go mod tidy && go mod download && go clean --modcache

RUN cd services/user && go build -o twitter-clone-user ./cmd/main.go

EXPOSE 7003

CMD [ "services/user/twitter-clone-user" ]