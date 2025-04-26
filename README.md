# Distributed Database in Go

## Overview

This is a basic distributed database system with a master node and slave nodes, using Go and SQLite.

## Features

- SQLite as the storage engine
- Master replicates operations to slaves
- HTTP communication between nodes
- Simple schema and record operations

## How to Run

1. Install dependencies:

```bash
go mod tidy
```

2. Run slave node:

```bash
go run node/main.go
```

3. Run master node (in another terminal):

```bash
go run master/main.go
```

## Author

Pavly
