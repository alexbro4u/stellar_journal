FROM golang:latest
WORKDIR go/src/app
COPY . .
WORKDIR ./cmd/stellar_journal

RUN go mod download
RUN GOOS=linux go build main.go

RUN chmod +x main
CMD ["./main"]
