FROM golang:1.22

ADD go.mod go.sum /

RUN go mod download

WORKDIR /app

COPY . .

ENV GO_ENV=production

RUN go build -o /app/visa_appointment_watcher cmd/visa_appointment_watcher/main.go

EXPOSE 3000
CMD ["/app/visa_appointment_watcher"]

