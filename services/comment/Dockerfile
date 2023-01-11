FROM golang:latest

ENV GOPATH=/

RUN apt-get update
RUN apt-get -y install postgresql-client

COPY services/ services

RUN chmod +x services/comment/wait-for-postgres.sh

COPY go.work .

COPY services/comment/go.mod services/comment

RUN cd services/comment && go mod tidy && go mod download && go clean --modcache

RUN cd services/comment && go build -o twitter-clone-comment ./cmd/main.go

EXPOSE 7005

CMD [ "./main" ]