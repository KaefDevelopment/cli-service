FROM golang:latest

WORKDIR /app

COPY . .

RUN go build -ldflags "-X main.version=$(git describe --tags --abbrev=0)" -o myapp cmd/cli/main.go

EXPOSE 8181

CMD ["./myapp"]


