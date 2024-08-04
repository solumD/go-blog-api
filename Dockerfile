FROM golang:1.22

WORKDIR /usr/local/src/application

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY ./ ./
RUN go build -o ./bin/app.exe cmd/app/main.go

CMD ["./bin/app"]