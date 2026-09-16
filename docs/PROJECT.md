# Pulse

Real-time event-driven notification system (learning project).

Core eventual flow:

```text
Alice ──follow──► Go API ──► PostgreSQL (Follow / Notification / Outbox)
                                    │
                                    ▼
                              Event system
                           ┌───────┴───────┐
                           ▼               ▼
                       WebSocket         Email
                           │
                           ▼
                          Bob
```

Collaboration model, phases, and agent rules live in `.cursor/rules/`.

References (introduced when relevant):

- [Transactional outbox](https://microservices.io/patterns/data/transactional-outbox)
- [Idempotent consumer](https://microservices.io/patterns/communication-style/idempotent-consumer.html)
