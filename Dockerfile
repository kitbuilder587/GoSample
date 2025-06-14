FROM golang:1.24.3

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

# Создаём папку вручную
RUN mkdir -p bin && go build -o bin/cryptotrack ./cmd/server
RUN ls -l /app/docs


CMD ["./bin/cryptotrack"]
