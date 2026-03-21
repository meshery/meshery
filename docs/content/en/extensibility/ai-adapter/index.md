---
title: "Extensibility: AI Adapter"
description: Natural language to infrastructure design using pluggable LLM providers.
---

## What is the AI Adapter?

The Meshery AI Adapter lets you describe infrastructure in plain English and get a Kubernetes manifest back automatically.

**Example:**
```
"Deploy nginx with 2 replicas on Kubernetes"
        ↓
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  replicas: 2
  ...
```

---

## Bring Your Own Model (BYOM)

You can use any LLM — local or cloud — without changing any code.

| Provider  | Type            | API Key needed |
|-----------|-----------------|----------------|
| Ollama    | Local (private) | No             |
| OpenAI    | Cloud           | Yes            |
| Anthropic | Cloud           | Yes            |

---

## How to use

### With Ollama (local, free, private)

```bash
# 1. Install and start Ollama
ollama pull llama3
ollama serve

# 2. Send a request
curl -X POST http://localhost:9081/api/ai/generate \
  -H "Content-Type: application/json" \
  -d '{
    "prompt": "Deploy nginx on Kubernetes with 2 replicas",
    "provider": "ollama"
  }'
```

### With OpenAI

```bash
export OPENAI_API_KEY="sk-..."

curl -X POST http://localhost:9081/api/ai/generate \
  -d '{"prompt": "Deploy nginx", "provider": "openai"}'
```

---

## API Endpoints

### POST `/api/ai/generate`
Generate a design from a natural language prompt.

**Request:**
```json
{
  "prompt": "Deploy nginx with 2 replicas",
  "provider": "ollama"
}
```

**Response:**
```json
{
  "provider": "ollama/llama3",
  "description": "A 2-replica nginx deployment.",
  "design": "apiVersion: apps/v1\nkind: Deployment..."
}
```

### GET `/api/ai/providers`
List available providers.

```json
{
  "providers": [
    { "id": "ollama", "available": true },
    { "id": "openai", "available": false }
  ]
}
```

---

## Running Tests

```bash
# Unit tests (no Ollama needed)
go test ./server/handlers/ -run "^Test[^I]" -v

# Integration tests (Ollama must be running)
go test ./server/handlers/ -run Integration -v -timeout 300s
```

---

## Related

- Issue: [#17097](https://github.com/meshery/meshery/issues/17097)