# promethea

## Quickstart 🚀

Run locally (fast):

```bash
make run-ollama  # requires locally installed ollama (e.g. via brew)
make run-app

make error-missing-state
make error-invalid-country
```

Run containerized (slow):

```bash
make up  # run containers
make pull  # pull image in the ollama container

make generate  # ask the model a question
```
