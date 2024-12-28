# Abyssal Kraken

Abyssal Kraken is a library to implement **Event Sourcing**, **Domain Events**, and **Aggregate Roots** patterns in Go. It is designed to be simple and extensible, making it easy to manage events and maintain consistency in the Domain Model.

---

## 🛠️ Installation

Add the library to your project using **Go Modules**:

```bash
go get github.com/abyssal-kraken/abyssalkraken
```

---

## ✨ Key Features

- **Aggregate Root**:
    - Generic interface for aggregate modeling.
    - Simple implementation through `SimpleAggregateRoot`.

- **Domain Event**:
    - Interface for domain events.
    - Tracking of pending events for processing.

- **Event Handlers**:
    - Support for transactional and non-transactional handlers.
    - Concurrent execution for non-transactional handlers.

---

## 🚀 How to Use

Here’s a basic example of how to use the library:

### Aggregate with Domain Events

```go
package main

import (
	"fmt"
	"time"

	"github.com/abyssal-kraken/abyssalkraken"
)

type MyAggregateID string

func (id MyAggregateID) String() string {
	return string(id)
}

type MyEvent struct {
	ID      abyssalkraken.EventID
	Msg     string
	Created time.Time
}

func (e *MyEvent) AggregateID() abyssalkraken.AggregateID {
	return MyAggregateID(e.ID.String())
}

func (e *MyEvent) EventID() abyssalkraken.EventID {
	return e.ID
}

// etc...

func main() {
	agg := abyssalkraken.NewSimpleAggregateRoot[MyAggregateID, *MyEvent](MyAggregateID("aggregate-1"))
	event := &MyEvent{
		ID:      abyssalkraken.RandomEventID(),
		Msg:     "Event Triggered",
		Created: time.Now(),
	}

	// Add event to the aggregate
	agg.AddEvent(event)

	fmt.Printf("Pending events: %v\n", agg.HasPendingEvents())
	fmt.Printf("Collected events: %v\n", agg.CollectPendingEvents())
}
```

---

## 🌟 Advanced Features

### Working with Handlers

The library supports handlers and transactional execution for domain events. Here's an example:

```go
type MyEventHandler struct {}

func (h *MyEventHandler) Transactional() bool {
    return false
}

func (h *MyEventHandler) Handle(aggregateRoot *abyssalkraken.SimpleAggregateRoot[MyAggregateID, *MyEvent], domainEvent *MyEvent) error {
    fmt.Println("Handling event:", domainEvent.Msg)
    return nil
}
```

---

## 🧪 Testing

Basic examples are included in the `examples/` directory. For more details, check the file. If you'd like to run the unit tests directly:

```bash
go test ./...
```

---

## 📂 Project Structure

.
├── abyssalkraken                 # Core source code of the library
│   ├── aggregate_root.go         # Base interface for aggregates
│   ├── domain_event.go           # Base interface for domain events
│   ├── domain_event_handler.go   # Handlers for domain events
│   ├── domain_event_publisher.go # Interface for event publishing
│   ├── domain_event_store.go     # Simple event storage
│   └── simple_aggregate_root.go  # Basic implementation of Aggregate Root
├── go.mod                        # Go module definition
└── README.md                     # Project documentation

---

## 🤝 Contributing

Contributions are welcome! For more information, check the [CONTRIBUTING.md](./CONTRIBUTING.md) file.

---

## 🛡️ License

This project is licensed under the [MIT License](./LICENSE).
