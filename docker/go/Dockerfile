FROM golang:1.20.4

COPY ./go.mod /src/go.mod
WORKDIR /src
RUN curl -s https://packagecloud.io/install/repositories/golang-migrate/migrate/script.deb.sh | bash && apt-get update && apt-get install migrate jq
RUN go mod download

COPY . /src

# CMD go run main.go