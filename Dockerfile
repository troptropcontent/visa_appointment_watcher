FROM golang:1.22

ADD go.mod go.sum /

RUN go mod download

WORKDIR /app

COPY . .

ENV GO_ENV=production

# Install Tailwind
RUN curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-arm64
RUN chmod +x tailwindcss-linux-arm64
RUN mv tailwindcss-linux-arm64 ./bin/tailwindcss
# Build the Tailwind CSS file
RUN ./bin/tailwindcss -i ./internal/assets/stylesheets/tailwind_input.css -o ./public/css/tailwind_output.css --minify

RUN go build -o /app/visa_appointment_watcher cmd/visa_appointment_watcher/main.go

EXPOSE 1234
CMD ["/app/visa_appointment_watcher"]

