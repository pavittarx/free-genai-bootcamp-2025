## Opea Comps
The goal is to run the Ollama API in a container.
Model: llama3.2:1b

Technical Requirements
- Docker containers will be used.
- The model will be llama3.2:1b
- The model will mount within the container.

Technical Uncertainity
- Not sure about the enviornment variables and how to set them.
- Whether or not my system will be able to run the required model.
- How to know if the model is up and running?
- The env variables, what to set and what not?
- How to use the model once it is downloaded and running?

Learnings / Observations / Improvements
- The model needs to be pulled after the container is up and running.
- The model is pulled by making an API request to the Ollama server.
- The response of the model is in JSON with each word. 
- The API needs to be triggered from the host machine. The docker container port is mapped to the host machine port. 
- The host machine port is in most cases 4 digits long. The docker container uses longer port numbers.
- The docker container port is mapped to the host machine port. 
- The model is saved within the container, so it will vanish once the container is done. We need to mount it to local file system to avoid that. 

### Get the HOST IP
```sh
    hostname -I | awk '{print $1}'
```

## Run Server
```sh
    HOST_IP=$(hostname -I | awk '{print $1}') NO_PROXY=localhostLLM_ENDPOINT_PORT=8008
    LLM_MODEL_ID="llama3.2:1b"
    docker compose up
```

### Docs
- Ollama: https://github.com/ollama/ollama/blob/main/docs/api.md
- Docker:https://docs.docker.com/engine/install/ubuntu/
- Docker Extension on VSCode

## Model API requests

### Pull the Model
```sh
curl http://localhost:8008/api/pull -d '{
  "model": "llama3.2:1b"
}'
```

### Generate Responses

```sh
curl http://localhost:8008/api/generate -d '{
  "model": "llama3.2:1b",
  "prompt": "Why is the sky blue?"
}'
```