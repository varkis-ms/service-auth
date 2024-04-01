FROM golang:1.21.3
ENV GIN_MODE=release
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go mod download
RUN go mod tidy
RUN go build ./cmd/auth
ENTRYPOINT /app/auth