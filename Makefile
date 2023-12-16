# --- Run containerized ---

up: down
	docker-compose up --build

down:
	docker compose down --remove-orphans

app-restart:
	docker compose down --remove-orphans app && docker compose up --build app


# --- Run natively ---

run-ollama:
	ollama serve

run-app:
	go run cmd/rest-server/main.go


# --- API calls ---

pull:
	curl -X POST -H "Content-Type: application/json" -d '{"name":"mistral"}' http://127.0.0.1:8080/api/pull

generate:
	curl -X POST -H "Content-Type: application/json" -d '{"model":"mistral", "stream":false, "prompt":"what is 1+1?"}' http://127.0.0.1:8080/api/generate

error-missing-state:
	curl -X POST -H "Content-Type: application/json" -d '{"message":"Missing value for field \"state\"", "code":"412", "value":""}' http://127.0.0.1:8080/api/errors/guess_field


error-invalid-country:
	curl -X POST -H "Content-Type: application/json" -d '{"message":"not a valid country code", "code":"", "value":"GE"}' http://127.0.0.1:8080/api/errors/guess_field
