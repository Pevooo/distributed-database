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

## Project Structure
```text
distributed-database/
|----- send_dummy.go
|----- show_all/main.go
|----- node/main.go
|----- master/main.go
|----- README.md
|----- go.sum
|----- go.mod
```

## Files

1. `master/main.go`: The master file.

2. `node/main.go`: The slave file.

3. `show_all/main.go`: A helper file that may be used to view the databases of the master and slaves. Note that this file is not a part of the project. It's just a helper. You can easily run the file by executing the following command from the root directory of the project:
```bash
go run show_all/main.go <num_databases>
```

`num_databases` is basically the `number of slaves + 1` which indicates the number of slaves. Note that each node (master or slave) has its own replica database so the total number of databases is `number of slaves + 1`.

4. `send_dummy.go`: A helper file that we may use to send a post request without using postman. Note that it's not a part of the project, it's just a helper file. We recommend using postman

## Logic

We define privileged operations as any operation that is related to creating or deleting of databases or tables

### On Creating or Deleting a database
This feature is not supported as we use sqlite, so we just bind to a database

### On Creating or Deleting a table
This is considered a privileged operation, so only the master can process this request and order the slaves to replicate it

### On running a simple query on master
The query will be run on the master the replicated to the slaves ensuring consistent replicas of the databases at all time

### On running a simple query on a slave
The query will be run on this slave and replicated to the master and the rest of the slaves

*Note that the replicated requests have a specific flag indicating that this request is replicated so that no extra replication happends*

## Technology
1. sqlite for the databse engine
2. GoLang for the code and logic
3. HTTP for communication between the master, slaves and external clients

## Authors

1. Pavly Samuel
2. John Ashraf
3. Ahmed Aziz
4. Abdelrahman Ayman
5. Abdelrahman Abdelhameed
