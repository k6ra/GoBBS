FROM golang:1.19

WORKDIR /opt

COPY . .
RUN go mod download
RUN go build -o app cmd/gobbs/main.go

CMD ["/opt/app"]