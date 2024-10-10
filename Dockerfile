ARG GO_VERSION=1
FROM golang:${GO_VERSION}-bookworm as builder

WORKDIR /usr/src/app
COPY go.mod ./
RUN go mod download && go mod verify
COPY . .
RUN go build -v -o /run-app .


FROM debian:bookworm
RUN apt-get update
RUN apt-get install openssh-client -y
COPY --from=builder /run-app /usr/local/bin/
CMD ["run-app"]
