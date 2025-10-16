FROM golang:1.24-alpine AS builder

RUN go install github.com/air-verse/air@v1.61.5

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

FROM golang:1.24-alpine AS development

RUN go install github.com/air-verse/air@v1.61.5

RUN apk add --no-cache git

WORKDIR /app

COPY . .

RUN if [ ! -f .air.toml ]; then air init; fi

EXPOSE 8080

CMD ["air", "-c", ".air.toml"]
