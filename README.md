# promethea

## Quickstart 🚀

Define prompt and valid data for the generated payload in the `instructions` folder.

### Locally/native (fast)

```bash
make run-ollama  # requires locally installed ollama (e.g. via brew)
make run-app
make pull-mistral  # pull the mistral LLM

# example requests
make error-missing-state
make error-invalid-country
```

Example results:

```bash
$ time make error-missing-state
curl -X POST -H "Content-Type: application/json" -d '{"filename":"errors.yaml", "message":"Missing value for field \"state\"", "code":"412", "value":""}' http://127.0.0.1:8080/api/errors/guess_field

{
"known_error_id": "state",
"description": "The field was missing and must contain a valid value.",
"original_error_message": "Missing value for field \"state\"",
"original_error_code": "412",
"original_error_value": ""
}

0.01s user 0.01s system 0% cpu 3.905 total
```

```bash
$ time make error-invalid-country
curl -X POST -H "Content-Type: application/json" -d '{"filename":"errors.yaml", "message":"not a valid country code", "code":"", "value":"GE"}' http://127.0.0.1:8080/api/errors/guess_field

{
"known_error_id": "country",
"description": "The country code provided was invalid.",
"original_error_message": "not a valid country code",
"original_error_code": "",
"original_error_value": "GE"
}

0.01s user 0.01s system 0% cpu 4.186 total
```

### Containerized (slow)

```bash
make up  # run containers
make pull  # pull the mistral LLM in the ollama container

# example request
make generate
```
