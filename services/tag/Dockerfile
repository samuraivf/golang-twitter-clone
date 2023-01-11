FROM golang:latest

ENV GOPATH=/

RUN apt-get update
RUN apt-get -y install postgresql-client

COPY services/ services

RUN chmod +x services/tag/wait-for-postgres.sh

COPY go.work .

COPY services/tag/go.mod services/tag

RUN cd services/tag && go mod tidy && go mod download && go clean --modcache

RUN cd services/tag && go build -o twitter-clone-tag ./cmd/main.go

EXPOSE 7007

CMD [ "./main" ]