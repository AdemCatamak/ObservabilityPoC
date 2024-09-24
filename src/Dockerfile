FROM golang:1.22

WORKDIR /app

COPY . /app/
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-gs-ping

CMD ["/docker-gs-ping"]