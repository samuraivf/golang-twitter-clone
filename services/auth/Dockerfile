FROM golang:latest

ENV GOPATH=/

COPY ./ ./

RUN apt-get update

COPY services/ services

COPY go.work .

COPY services/auth/go.mod services/auth
COPY services/auth/go.sum services/auth

RUN cd services/auth && go mod tidy && go mod download && go clean --modcache

RUN cd services/auth && go build -o twitter-clone-auth ./cmd/main.go

EXPOSE 7001

CMD [ "./main" ]