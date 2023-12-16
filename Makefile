# --- Container management ---

up: down
	docker-compose up --build # -d

down:
	docker compose down --remove-orphans

app-restart:
	docker compose down --remove-orphans app && docker compose up --build app


# --- Run app natively ---

app-native:
	go run cmd/rest-server/main.go


# --- API calls ---

pull:
	curl -X POST -H "Content-Type: application/json" -d '{"name":"mistral"}' http://127.0.0.1:8080/api/pull

generate:
	curl -X POST -H "Content-Type: application/json" -d '{"model":"mistral", "stream":false, "prompt":"what is 1+1?"}' http://127.0.0.1:8080/api/generate
