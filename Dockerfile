FROM golang:1.22.2

WORKDIR /usr/src/app

COPY . .

RUN go mod download 

ENV TODO_PORT=7540 \
    TODO_DBFILE="/usr/src/app/scheduler.db" \
    TODO_WEB_DIR="/usr/src/app/web" \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

RUN  go build -o /usr/src/app/cmd/todo/main ./cmd/todo/main.go

CMD ["/usr/src/app/cmd/todo/main"] 