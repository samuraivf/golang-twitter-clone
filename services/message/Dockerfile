FROM golang:latest

ENV GOPATH=/

RUN apt-get update
RUN apt-get -y install postgresql-client

COPY services/ services

RUN chmod +x services/message/wait-for-postgres.sh

COPY go.work .

COPY services/message/go.mod services/message

RUN cd services/message && go mod tidy && go mod download && go clean --modcache

RUN cd services/message && go build -o twitter-clone-message ./cmd/main.go

EXPOSE 7004

CMD [ "./main" ]