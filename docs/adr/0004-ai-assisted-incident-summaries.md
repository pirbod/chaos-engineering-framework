# ADR 0004: AI-Assisted Incident Summaries

## Status

Accepted

## Context

Teams need fast, consistent summaries during experiments, but local demos should not require paid AI services or secrets.

## Decision

Implement a provider interface with a deterministic local provider by default and OpenAI/Azure OpenAI-compatible configuration placeholders.

## Consequences

The demo is reliable offline, and production teams can later plug in an approved provider without changing API contracts.
