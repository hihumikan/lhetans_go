FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ./ ./
RUN go build -o lhetansgo

FROM golang:alpine AS runner

WORKDIR /app

COPY --from=builder /app/lhetansgo ./

EXPOSE 3000

CMD [ "/app/lhetansgo" ]
