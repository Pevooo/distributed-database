
# 🗄️ Distributed Database in Go 🚀

## 📘 Overview

Welcome to **Distributed Database in Go**, a lightweight yet powerful project that simulates a distributed database system using **Go** and **SQLite**. This system is built with **replication** in mind—ensuring that every node (master or slave) maintains a **consistent copy** of the database. The nodes communicate through **HTTP POST** requests, making the system simple to interact with and easy to extend.

### 🔑 Key Concepts

- 🧠 The **master node** processes **privileged operations** such as `CREATE TABLE`, and replicates them to the slaves.
- 🧑‍🤝‍🧑 **Slave nodes** handle regular queries but cannot execute database or table creation commands.
- 🧬 Every node (master or slave) maintains its own **replica** of the SQLite database.
- 🔄 All nodes share a consistent state by replicating commands to each other.

---

## ✨ Features

- 🛢️ SQLite as the storage engine
- 📡 HTTP-based communication between nodes
- 👑 Master-slave architecture with replication
- 🧪 Support for SQL operations (excluding DB creation)
- 🧱 Schema and record-level operations
- 🔐 Control over privileged queries via master authority

---

## 🔧 Prerequisites

- ✅ [Go 1.16+](https://golang.org/dl/)

---

## 🚀 How to Run

### 1. Install dependencies

```bash
go mod tidy
```

### 2. Run Slave Nodes

Each slave requires a unique ID (starting from 1) and the total number of slaves:

```bash
go run node/main.go <slave_id> <num_slaves>
```

📌 Example:
```bash
go run node/main.go 1 2
```

### 3. Run Master Node

In a separate terminal window:

```bash
go run master/main.go <num_slaves>
```

📌 Example:
```bash
go run master/main.go 2
```

---

## 🧪 How to Use

Use Postman or any other API client to interact with your distributed database via HTTP POST requests. Every node hosts its own SQLite database that reflects all replicated operations.

### ✅ Example: Create a Table

Send this JSON to the **master node** (usually on `http://localhost:8080`):

```json
{
  "command": "exec",
  "query": "CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, username TEXT UNIQUE, email TEXT, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP)",
  "is_master": false
}
```

⚠️ **Important:** This request will only work when sent to the **master**. If sent to a **slave**, it will be rejected as a privileged operation.

✅ Once the master processes the request, it sends a **privileged replicated** version to all slaves so that their databases reflect the same schema.

---

## 📁 Project Structure

```text
distributed-database/
├── send_dummy.go           # Helper: Sends test POST requests
├── show_all/main.go        # Helper: Shows all node databases
├── node/main.go            # Slave node logic
├── master/main.go          # Master node logic
├── README.md
├── go.sum
└── go.mod
```
## Project Archeticture

+-------------------------------------------------------+
|                      Master Node                      |
|                                                       |
| - Privileged Instructions                            |
| - Broadcast to Slaves                                |
|                                                       |
+-------------------------------------------------------+
                            |
                            |
              +-------------+-------------+
              |                           |
              |                           |
+-------------+-------------+   +---------+-------------+
|         Slave Node        |   |       Slave Node       |
|                           |   |                        |
| - Standard Access DB      |   | - Standard Access DB   |
| - Broadcast to All        |   | - Broadcast to All     |
|                           |   |                        |
+---------------------------+   +------------------------+
---

## 📄 File Descriptions

- `master/main.go`: Contains logic for the master node—handles privileged queries and replication.
- `node/main.go`: Logic for slave nodes—executes standard queries and receives replication from the master.
- `show_all/main.go`: Utility to print all records from all nodes (for debugging/testing).
  ```bash
  go run show_all/main.go <num_databases>
  ```
  Replace `<num_databases>` with `num_slaves + 1`.

- `send_dummy.go`: Sends HTTP POST requests programmatically (alternative to Postman).

---

## 🧠 System Logic

### 🔐 Privileged Operations

Operations like `CREATE TABLE` are **privileged** and only allowed on the **master node**. Once processed, the master replicates these commands to slaves.

### 🔄 Query Replication Logic

- 🧑‍💻 **Query sent to master** → executed → replicated to slaves.
- 🧑‍🔧 **Query sent to slave** → executed → replicated to master and other slaves.

✔️ A flag is used in the JSON body to detect **replicated requests**, preventing infinite loops.

---

## ⚙️ Technologies Used

- 🗃️ **SQLite** – Embedded database engine.
- 🧬 **Go (Golang)** – Main programming language for logic and HTTP handling.
- 🌐 **HTTP** – Communication between nodes and external clients.

---

## 👨‍💻 Authors

- Pavly Samuel
- John Ashraf
- Ahmed Aziz
- Abdelrahman Ayman
- Abdelrahman Abdelhameed

---

## ❤️ Acknowledgments

Thanks to the entire team for building and documenting this educational distributed system project! 🙌

---
