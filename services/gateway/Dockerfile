FROM golang:latest

ENV GOPATH=/

COPY ./ ./

RUN apt-get update

COPY services/ services

COPY go.work .

COPY services/gateway/go.mod services/gateway
COPY services/gateway/go.sum services/gateway

RUN cd services/gateway && go mod tidy && go mod download && go clean --modcache

RUN cd services/gateway && go build -o twitter-clone-gateway ./cmd/main.go

EXPOSE 7000

CMD [ "services/gateway/twitter-clone-gateway" ]