server:
	go run cmd/visa_appointment_watcher/main.go
tailwind:
	./bin/tailwindcss -i ./internal/assets/stylesheets/tailwind_input.css -o ./public/css/tailwind_output.css --watch

