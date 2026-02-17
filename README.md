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

These wrap the standard library’s `sync.Mutex` and `sync.RWMutex`. The main addition is a consistent API that includes **TryLock** (and **TryRLock** for RWMutex), so callers can attempt to acquire the lock without blocking and react if it’s already held.

### Condition variable (Cond)

**Cond** wraps `sync.Cond`. It is always used with an external lock (the caller passes a `sync.Locker` to `New`). **Wait** atomically releases the lock and blocks the goroutine until **Signal** (one waiter) or **Broadcast** (all waiters) is called; when **Wait** returns, the lock is held again. The typical pattern is: hold the lock, check a condition in a loop, call **Wait** inside the loop if the condition isn’t met, then recheck after **Wait** returns. This avoids lost wakeups and spurious wakeups.

### Barrier

The barrier is implemented with a **mutex and a condition variable**. A shared counter counts how many goroutines have called **Wait**; when the count reaches N, the last goroutine resets the counter, increments a generation number, and calls **Broadcast** on the condition variable so all waiters are released. The generation number is used so that waiters from the previous round don’t wake up in the next round. Each goroutine that isn’t the Nth one waits on the condition until the generation changes, then returns. The barrier is reusable: after each “round,” all N goroutines can call **Wait** again for the next round.

## Usage and examples

See the package doc comments and the `*_test.go` files in each package for usage and examples.

## License

MIT
