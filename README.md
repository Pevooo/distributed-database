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

## How to Use

You can easily use the project by sending a post request to a slave or the master (each node has a replica of the database). We recommend using Postman for its user-friendly interface.

We can send a simple json to the master to create a table:
```json
{
  "command": "exec",
  "query": "CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, username TEXT UNIQUE, email TEXT, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP)",
  "is_master": false
}
```

This command will create a table with an `id` as a primary key, a unique text username, and a `created_at` tmiestamp with the current time. Note that this command will fail if we try to send the same json to a slave instead of the master. We can easily send this json to the master by sending it to `localhost:8080` which is the address of the master. After executing this request, the master will automatically send a previllaged replicated request to the slaves indicating that they have to create the same table with the same attributes and constraints.

## Files

3. `master/main.go`: The master file.

2. `node/main.go`: The slave file.

1. `show_all/main.go`: A helper file that may be used to view the databases of the master and slaves. Note that this file is not a part of the project. It's just a helper. You can easily run the file by executing the following command from the root directory of the project:
```bash
go run show_all/main.go <num_databases>
```

`num_databases` is basically the `number of slaves + 1` which indicates the number of slaves. Note that each node (master or slave) has its own replica database so the total number of databases is `number of slaves + 1`.

4. `send_dummy.go`: A helper file that we may use to send a post request without using postman. Note that it's not a part of the project, it's just a helper file. We recommend using postman
