# 0001: Start with an executable, dependency-free core

Status: accepted for the sample

We want the environment to be understandable and verifiable on first use. A small Go service demonstrates business logic, HTTP contracts, architecture checks, and CI without accounts or installation.

We use net/http and testing rather than adding chi and testify to this example. No frontend, persistence, Docker, or live model integration is needed for the exercise. This deliberately narrows the workspace's normal stack for a sample; a real product can add the standard stack with the usual dependency approval.

Consequences: quick onboarding and deterministic checks, but this is not a full product starter or a sandbox implementation. Runtime isolation, hosting protection, and integrations are adoption work described in docs/environment.md. Revisit when a product requirement makes those components necessary.
