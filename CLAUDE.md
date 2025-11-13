# CLAUDE.md - Development Guidelines

Development guidelines and best practices for this project.

For node: use 'bun' and 'typescript'

For python: use 'uv' and 'pydantic'

---

## SOLID Principles

Five design principles that make software more maintainable, flexible, and scalable.

### Single Responsibility (SRP)
Each class should have only one reason to change, with one specific responsibility.
- Separate UI widgets from business logic
- Keep repositories focused on data operations only
- Isolate validation logic into dedicated validator classes
- Benefits: easier testing, clearer code purpose, simpler maintenance

### Open/Closed (OCP)
Software entities should be open for extension but closed for modification.
- Use abstract classes and interfaces to define contracts
- Extend functionality by creating new implementations, not modifying existing code
- Example: Create `PaymentMethod` interface, then extend with `CreditCard`, `PayPal`, etc.
- Benefits: reduces bugs in existing code, safer to add features

### Liskov Substitution (LSP)
Objects of a subclass must be substitutable for objects of their parent class.
- Subclasses should strengthen, not weaken, parent class behavior
- Don't throw exceptions in overridden methods that parent doesn't throw
- Example: If `Bird` has `move()`, all bird subclasses should implement valid movement
- Benefits: predictable behavior, safer inheritance hierarchies

### Interface Segregation (ISP)
Clients shouldn't be forced to depend on interfaces they don't use.
- Create small, focused interfaces instead of large, monolithic ones
- Split `Worker` interface into `Workable`, `Eatable`, `Sleepable`
- Classes implement only the interfaces they need
- Benefits: more flexible code, easier to implement and test

### Dependency Inversion (DIP)
Depend on abstractions, not concrete implementations.
- High-level modules shouldn't depend on low-level modules
- Use dependency injection to provide implementations
- Define abstract `DataSource`, inject `ApiClient` or `LocalDatabase`
- Benefits: easier testing with mocks, flexible architecture, decoupled code

---

## DRY Principle (Don't Repeat Yourself)

- Extract repeated UI patterns into reusable widgets
- Use Dart mixins to share functionality across classes
- Separate business logic from UI components
- Create utility functions for common operations
- Benefits: less code, easier maintenance, fewer bugs, better testing

---

## KISS Principle (Keep It Simple, Stupid)

- Use Flutter's built-in widgets instead of creating complex custom solutions
- Write self-explanatory code with clear variable/function names
- Avoid over-engineering simple problems
- Minimize external dependencies
- Break down complex widgets into smaller, manageable pieces
- Start simple, add complexity only when necessary

---

## YAGNI Principle (You Aren't Gonna Need It)

Don't implement functionality until it's actually needed.
- Resist the urge to build features "just in case" they might be useful later
- Focus on current requirements, not hypothetical future needs
- Don't create abstract layers for one implementation
- Avoid premature optimization before measuring performance
- Don't build configuration systems until you need configurability
- Wait for actual use cases before adding flexibility
- Benefits: less code to maintain, faster delivery, lower complexity, easier to change

---

## Summary

Following these principles results in:
- Maintainable, extendable code
- Fewer bugs and faster debugging
- Better team collaboration
- Professional quality standards

**Remember: Good code is simple, clear, and purposeful.**
