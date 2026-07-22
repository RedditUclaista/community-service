FROM golang:1.25.0-bookworm

RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    && rm -rf /var/lib/apt/lists/*

RUN go install github.com/air-verse/air@latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

EXPOSE 10005

CMD ["air", "-c", ".air.toml"]

