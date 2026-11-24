# AI Adapter

## Overview

The AI Adapter connects Meshery to Large Language Models (LLMs), enabling users to describe
infrastructure in plain language and receive generated designs and configurations.

This reduces the manual effort of constructing complex topologies and makes Meshery more
accessible to engineers across different experience levels.

## Key Idea

Instead of building infrastructure visually or through configuration files, users can describe
their intent in natural language:

> "Deploy a microservices application with an API gateway, three backend services, and a Redis cache on Kubernetes."

The AI Adapter processes this input and converts it into a Meshery-compatible design,
including components and their relationships.

## Architecture Flow

1. **User Input** — A natural language prompt is provided through the Meshery UI or API.
2. **AI Adapter** — The adapter routes the request to a configured LLM provider such as OpenAI or Ollama.
3. **LLM Processing** — The model interprets the request and returns structured infrastructure data.
4. **Design Generation** — The adapter converts the output into a Meshery design and topology.
5. **Review and Deploy** — The generated design appears in Meshery’s canvas for review and deployment.

## LLM Provider Support

The AI Adapter is provider-agnostic and supports multiple LLM backends, including OpenAI, Ollama,
and custom models through BYOM (Bring Your Own Model).

BYOM enables teams to integrate private or fine-tuned models, ensuring control over model behavior
and data privacy.

## Future Enhancements

- Context-aware suggestions based on existing infrastructure and past designs.
- Multi-turn conversation support for iterative design refinement.
- Policy-aware generation aligned with organizational standards.
- Expanded compatibility with additional open-source and enterprise LLM providers.
