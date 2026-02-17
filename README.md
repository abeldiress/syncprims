# syncprims

A Go library of synchronization primitives: semaphores, mutexes, condition variables, and barriers. Each primitive is implemented in its own package so you can depend only on what you need.

## Primitives

| Package    | Primitive   | Description |
|-----------|-------------|-------------|
| `semaphore` | Semaphore  | Counting semaphore to limit concurrent access |
| `mutex`     | Mutex      | Mutual exclusion with `TryLock` |
| `rwmutex`   | RWMutex    | Read/write lock for reader-heavy workloads |
| `cond`      | Cond       | Condition variable for waiting on events |
| `barrier`   | Barrier    | Rendezvous for N goroutines |

## Installation

```bash
go get github.com/abeld/syncprims
```

Import the subpackages you need, for example:

```go
import (
    "github.com/abeld/syncprims/semaphore"
    "github.com/abeld/syncprims/mutex"
)
```

## High-level implementation overview

### Semaphore

The semaphore is implemented as a **buffered channel of empty structs**. The capacity of the channel is the number of permits. **Acquire** sends a value into the channel (blocking if the buffer is full), and **Release** receives from the channel (freeing one permit). This gives a counting semaphore with FIFO-style blocking and no extra state. **TryAcquire** uses a non-blocking send (`select` with `default`), and **AcquireContext** uses `select` between the send and `ctx.Done()` so acquisition can be cancelled or time-limited.

### Mutex and RWMutex

Both locks are implemented **purely with channels**, not `sync.Mutex`/`sync.RWMutex`.  
`Mutex` is a binary semaphore (`chan struct{}` of size 1) with `Lock`, `Unlock`, and `TryLock`.  
`RWMutex` combines a writer semaphore with a small internal mutex and reader count, allowing multiple readers or a single writer, plus non-blocking `TryLock`/`TryRLock`.

### Condition variable (Cond)

`Cond` is built on top of this library’s `Mutex`. Each waiter registers a dedicated channel and then unlocks and blocks. `Signal` takes one waiter from a slice and wakes it; `Broadcast` wakes them all. When a waiter is released it re-locks the associated mutex before returning.

### Barrier

The barrier uses the custom `Mutex` and `Cond` types. A counter tracks how many goroutines have called `Wait`; when it reaches `n`, the last goroutine resets the counter, bumps a generation counter, and broadcasts on the condition variable so all waiters continue. The generation counter allows the same barrier instance to be reused for multiple rounds.

## Usage and examples

See the package doc comments and the `*_test.go` files in each package for usage and examples.

## License

MIT
