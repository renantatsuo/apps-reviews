# gopq: Go Persistent-Queue

A lightweight, sqlite-backed persistent queue implementation in Go.

[![Go](https://github.com/mattdeak/gopq/actions/workflows/go.yml/badge.svg)](https://github.com/mattdeak/gopq/actions/workflows/go.yml)
[![GoDoc](https://godoc.org/github.com/mattdeak/gopq?status.svg)](https://godoc.org/github.com/mattdeak/gopq)

gopq is a lightweight, persistent queue implementation in Go, using SQLite as the underlying storage mechanism. It provides various queue types to suit different use cases, including simple queues, acknowledged queues, and unique item queues.

## Features

- Persistent storage using SQLite
- Unique/Non-Unique Queues
- Acknowledged/Non-Acknowledged Queues
- Blocking and non-blocking dequeue operations
- Context support for cancellation and timeouts
- Thread-safe operations
- Easy and Composable Dead Letter Queues

## Installation

To use gopq in your Go project, run:

```
go get github.com/mattdeak/gopq
```

Make sure you have SQLite installed on your system.

## Getting Started
Here's a minimal example to get you started with gopq:
```go
package main

import (
    "fmt"
    "github.com/mattdeak/gopq"
)

func main() {
    // Create a new simple queue
    queue, err := gopq.NewSimpleQueue("myqueue.db")
    if err != nil {
        panic(err)
    }
    defer queue.Close()

    // Enqueue an item
    err = queue.Enqueue([]byte("Hello, gopq!"))
    if err != nil {
        panic(err)
    }

    // Dequeue an item
    msg, err := queue.Dequeue()
    if err != nil {
        panic(err)
    }

    fmt.Println(string(msg.Item)) // Output: Hello, gopq!
}
```

## Usage

Here are some examples of how to use different queue types:

### Basic Queue
A simple FIFO queue.
```go
import "github.com/mattdeak/gopq"

// Create a new simple queue
queue, err := gopq.NewSimpleQueue("queue.db")
if err != nil {
    // Handle error
}

// Enqueue an item
err = queue.Enqueue([]byte("Hello, World!"))
if err != nil {
    // Handle error
}

// Dequeue an item (blocks until a message is available)
msg, err := queue.Dequeue()
if err != nil {
    // Handle error
}

fmt.Println(string(msg.Item))
```

### Acknowledged Queue
When an item is dequeued from an acknowledged queue, an ack deadline is set internally. If the ack deadline passes without an ack being received, the item becomes available at the front of the queue.

```go
import "github.com/mattdeak/gopq"
// Create a new acknowledged queue with custom options
queue, err := gopq.NewAckQueue("ack_queue.db", gopq.AckOpts{
    AckTimeout: 30 * time.Second,
    MaxRetries: 3,
    RetryBackoff: 5 * time.Second,
})
if err != nil {
// Handle error
}
// Enqueue an item
err = queue.Enqueue([]byte("Process me"))
if err != nil {
// Handle error
}
// Dequeue an item
msg, err := queue.Dequeue()
if err != nil {
// Handle error
}
// Process the item...
// Acknowledge the item if processing was successful
err = queue.Ack(msg.ID)
if err != nil {
// Handle error
}
// Or, if processing failed, negative acknowledge the item
err = queue.Nack(msg.ID)
if err != nil {
// Handle error
}
```

### Unique Queue
A unique queue ensures that every item in the queue is unique. Once an item is
successfully dequeued, it can be enqueued again immediately.

```go
import "github.com/mattdeak/gopq"

// Create a new unique queue
queue, err := gopq.NewUniqueQueue("unique_queue.db")
if err != nil {
    // Handle error
}

// Enqueue items (duplicates will be ignored)
queue.Enqueue([]byte("unique_item_1"))
queue.Enqueue([]byte("unique_item_1")) // This will be ignored
queue.Enqueue([]byte("unique_item_2"))

// Dequeue items
msg1, _ := queue.Dequeue()
msg2, _ := queue.Dequeue()
msg3, err := queue.TryDequeue() // This will return an error (queue is empty)

```

## Queue Methods

All queue types (SimpleQueue, AckQueue, UniqueQueue, UniqueAckQueue) provide the following methods for enqueueing and dequeueing:

### Enqueue Methods
* `Enqueue(item []byte) error`: Adds an item to the queue. Blocks if necessary.
* `EnqueueCtx(ctx context.Context, item []byte) error`: Adds an item to the queue with context support.
* `TryEnqueue(item []byte) error`: Attempts to add an item to the queue immediately, non-blocking.
* `TryEnqueueCtx(ctx context.Context, item []byte) error`: Attempts to add an item to the queue immediately with context support.

### Dequeue Methods
* `Dequeue() (Msg, error)`: Removes and returns an item from the queue. Blocks if the queue is empty.
* `DequeueCtx(ctx context.Context) (Msg, error)`: Removes and returns an item with context support.
* `TryDequeue() (Msg, error)`: Attempts to remove and return an item immediately, non-blocking.
* `TryDequeueCtx(ctx context.Context) (Msg, error)`: Attempts to remove and return an item immediately with context support.

All methods return a `Msg` struct containing the item and its ID, along with an error if the operation fails. For `AckableQueue`, the dequeue methods also update the acknowledgement deadline for the dequeued item.

### Method Patterns
- Methods without `Ctx` use a background context internally.
- `Try` methods are non-blocking and return immediately.
- Methods with `Ctx` allow for cancellation and timeouts using the provided context.

### Additional Methods for AckableQueue
* `Ack(id int64) error`: Acknowledges successful processing of an item.
* `Nack(id int64) error`: Indicates failed processing, potentially requeueing the item.

## Queue Types

1. **SimpleQueue**: Basic FIFO queue with no additional features. Ideal for simple task queues or message passing where order matters but acknowledgment isn't necessary.
2. **AckQueue**: Queue with acknowledgment support. Items must be acknowledged after processing, or they will be requeued after a timeout. Suitable for ensuring task completion in distributed systems or when processing reliability is crucial.
3. **UniqueQueue**: Queue that only allows unique items. Duplicate enqueue attempts are silently ignored. Useful for de-duplication scenarios or when you need to ensure only one instance of a task is queued.
4. **UniqueAckQueue**: Combination of UniqueQueue and AckQueue. Ensures unique items with acknowledgment support. Ideal for scenarios requiring both de-duplication and reliable processing.


## Advanced Features
### In-Memory Queue
If you don't require persistence, you can use an in-memory queue
by calling the constructor with an empty string.

```go
q := gopq.NewSimpleQueue("") // Now uses an in-memory database
```
This can be useful for testing.

### Dead Letter Queues and Failure Callbacks

GoPQ now supports dead letter queues through a more flexible callback system. Instead of directly specifying a dead letter queue, you can register failure callbacks that are called when a message fails to acknowledge after all retries have been exhausted.

To set up a dead letter queue:

1. Create an `AcknowledgeableQueue` for your main queue.
2. Create another queue to serve as your dead letter queue.
3. Use the `RegisterDeadLetterQueue` method to set up the dead letter functionality.

Example:

```go
mainQueue := NewAckQueue(...)
deadLetterQueue := NewSimpleQueue(...)
mainQueue.RegisterDeadLetterQueue(deadLetterQueue)
```

Dead letter queues can be any queue type or anything supporting the Enqueuer interface. You can use multiple queues with different settings, and even chain dead letter queues if needed. This method is a convenient shorthand for registering a failure callback that enqueues failed messages to the specified dead letter queue.

### Custom Failure Callbacks

For more complex scenarios, you can register custom failure callbacks:


```go

mainQueue.RegisterOnFailureCallback(func(msg Msg) error {
// Custom logic for handling failed messages
log.Printf("Message failed: %v", msg)
    return nil
})

```
See the examples directory for more.


### Configurable Retry Mechanism

AckQueue and UniqueAckQueue support configurable retry mechanisms:

```go
queue, := gopq.NewAckQueue("queue.db", gopq.AckOpts{
    AckTimeout: 1 * time.Minute, // Any message that takes longer than 1 minute to ack will be requeued.
    MaxRetries: 5 // 0 for no retries, -1 for infinite retries
    RetryBackoff: 10 * time.Second, // Sets a new ack deadline to be the max of (current deadline, now + retry backoff)
})
```

## Examples

There are several, more detailed examples demonstrating various features of gopq. These examples are located in the `examples` directory at the root of the project. To run an example, navigate to its directory and use `go run main.go`. For instance:

Available examples:
- Simple Queue Usage: `examples/simple_queue_example`
- Acknowledged Queue: `examples/ack_queue_example`
- Unique Queue: `examples/unique_queue_example`
- Dead Letter Queue: `examples/dead_letter_queue_example`
- Multiple Tiered Dead Letter Queues: `examples/tiered_dlq_example`
- Concurrent Queue Usage: `examples/concurrent_queue_example`

Each example demonstrates different features and use cases of gopq. We encourage you to explore these examples to better understand how to use the library in your projects.

## Contributing
Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change. Ensure to update tests as appropriate.


## License
This project is licensed under the MIT License.


### Future Work
- More Queue Configurability (LIFO, Priority Queues, etc.)
- Batch Enqueue/Dequeue operations
- Various Efficiency Improvements