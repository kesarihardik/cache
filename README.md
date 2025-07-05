# cachelib

**`cachelib`** is a generic, thread-safe in-memory caching library in Go. It provides pluggable eviction strategies via a clean `Cache` interface, with a built-in implementation of **LRU (Least Recently Used)**. Designed for flexibility, performance, and extensibility.

---

## Features

- ✅ **Generic** cache interface using Go 1.18+ type parameters
- 🔁 **LRU eviction** and **LFU eviction** built-in
- 🔒 **Thread-safe** with `sync.RWMutex`
- 🧩 Clean `Cache[K, V]` interface for extension
- ⚙️ Efficient in-memory design: `O(1)` `Get` and `Put`

---

## Installation

```bash
go get github.com/yourusername/cachelib
