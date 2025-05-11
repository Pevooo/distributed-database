# Distributed Database in Go

## Overview

This is a basic distributed database system with a master node and slave nodes, using Go and SQLite. Each node takes queries by sending a POST request to it. Slaves can't run database or table operations like `CREATE TABLE`, `CREATE DATABASE`, etc... . Each node (master or slave) has a separate database and sends the same query it gets to other nodes so that each database is a replica. In case of database or table queries, only the master node can process these queries and sends a previllaged query to the slaves so that the database it consistent across all nodes.

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

2. Run one or more slave nodes and give each slave a unique id starting from 1:

```bash
go run node/main.go <slave_id> <num_slaves>
```

3. Run master node (in another terminal):

```bash
go run master/main.go <num_slaves>
```
