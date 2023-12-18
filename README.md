# promethea

## Quickstart 🚀

Define prompt and valid data for the generated payload in the `instructions` folder.

Run locally (fast):

```bash
make run-ollama  # requires locally installed ollama (e.g. via brew)
make run-app

# example requests
make error-missing-state
make error-invalid-country
```

Run containerized (slow):

```bash
make up  # run containers
make pull  # pull image in the ollama container

# example request
make generate
```
