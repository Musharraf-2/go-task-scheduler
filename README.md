GOQ
===

GOQ is a self-hosted asynchronous task queue system written in Go.

It allows applications to execute background jobs reliably with concurrency, retries, scheduling, and observability.

Inspired by Sidekiq and Celery, GOQ focuses on production-grade reliability while remaining simple and extensible.

Features
--------

*   Concurrent worker pools
    
*   At-least-once delivery
    
*   Lease-based execution model
    
*   Automatic retries with backoff
    
*   Delayed & scheduled jobs
    
*   Dead-letter queue support
    
*   Crash-safe lease recovery
    
*   REST API
    
*   Web dashboard
    
*   Health & metrics endpoints
    
*   Memory and Postgres storage options
    

Use Cases
---------

*   Sending emails
    
*   Webhook processing
    
*   Report generation
    
*   Data processing
    
*   Background integrations
    
*   Scheduled jobs
    

Design Principles
-----------------

*   Framework-agnostic core
    
*   Explicit task state machine
    
*   Graceful shutdown
    
*   Observability-first design
    
*   Idempotent task support
    

Run Tests
---------

`   go test ./... -race   `

GOQ is a production-inspired Go project designed to demonstrate real-world concurrency, reliability, and backend systems engineering patterns.