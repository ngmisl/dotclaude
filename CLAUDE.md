# Development Guidelines

@AGENTS.md

**Stack:** Node (bun + typescript) | Python (uv + pydantic) | Go (typesafe + go.mod)

## SOLID Principles

**Single Responsibility (SRP)** - One class, one responsibility, one reason to change
- Separate concerns: UI, business logic, data access, validation

**Open/Closed (OCP)** - Open for extension, closed for modification
- Use interfaces/abstractions; extend via new implementations, not edits

**Liskov Substitution (LSP)** - Subclasses must be substitutable for parent classes
- Strengthen behavior, never weaken; maintain contracts

**Interface Segregation (ISP)** - Many specific interfaces > one general interface
- Clients implement only what they need

**Dependency Inversion (DIP)** - Depend on abstractions, not implementations
- Inject dependencies; decouple high-level from low-level modules

## Core Principles

**DRY (Don't Repeat Yourself)**
- Extract patterns into reusable components/functions/utilities

**KISS (Keep It Simple)**
- Prefer built-in solutions; clear names; avoid over-engineering; minimize dependencies

**YAGNI (You Aren't Gonna Need It)**
- Build only what's needed now; defer complexity until required

## Result
Code that is maintainable, testable, simple, clear, and purposeful.
