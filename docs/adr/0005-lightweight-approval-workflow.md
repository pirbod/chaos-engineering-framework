# ADR 0005: Lightweight Approval Workflow

## Status

Accepted

## Context

High and critical chaos experiments need explicit owner acknowledgement, rollback thresholds and evidence without turning the platform into a ticketing system.

## Decision

Add a small in-memory approval workflow to the Go API. Critical requests require runbook and rollback threshold metadata. Required approvers are derived from risk and blast radius.

## Consequences

The local demo stays database-free while showing how a production platform could add persistence, audit events and identity-aware approval checks.
