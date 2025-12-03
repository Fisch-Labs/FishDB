

<p align="center">
    <img src="fish.png" alt="FishDB Logo" width="120"/>
</p>

<h1 align="center">FishDB</h1>
<p align="center">
     Under Development · Lightweight, Embeddable, and Distributed Graph Database for AI
</p>

---

## Overview  

**FishDB** is a lightweight, embeddable, and distributed **graph database** written in Go, designed for high-performance **Retrieval-Augmented Generation (RAG)** applications.  
It aims to be the **fastest graph-naïve versioned database**, offering a powerful way to store, query, and manage evolving graph data.  

---

## Features  

- **Graph-based Data Model** – Store data as nodes and edges, with built-in versioning support.  
- **RAG-Optimized Performance** – Ultra-fast retrieval for AI & LLM-powered applications.  
- **Naive Versioning** – Efficiently manage and query multiple versions of graph data.  
- **Multiple Access Methods**  
    - RESTful API  
    - GraphQL API (with subscriptions)  
- **Clustering & Distribution** – Scale horizontally with multi-node clusters.  
- **Flexible Storage**  
    - In-memory mode for blazing speed  
    - Disk-based mode for persistence  
- **Interactive CLI Console** – Manage and query the database directly from the terminal.  

---

## Getting Started  

### Prerequisites  
- [Go](https://golang.org/doc/install) **1.12+**

### Build from Source  

```bash
# Clone repository
git clone https://github.com/Fisch-Labs/FishDB.git
cd FishDB

# Tidy dependencies
go mod tidy

# Build binary
go build ./cli/fishdb.go
```

### Run FishDB

**Disk-based storage**

```bash
./fishdb server -db-path /path/to/your/db
```

**In-memory mode**

```bash
./fishdb server -mem
```

**Note:** The server runs on HTTPS by default on port `9090`. The API endpoints are rooted at `/db`. For example, to access the info endpoint, use `https://localhost:9090/db/v1/info`.

---

## Usage

### REST API

Endpoints are available in the `api/v1` directory.

### GraphQL

Query and mutate graph data via `/db/graphql`.

### Interactive Console

```bash
./fishdb console
```

Manage users, run queries, and explore FishDB interactively.

---

## Contributing

Contributions are welcome!
- Open an issue to discuss new ideas or report bugs.
- Submit PRs for fixes or new features.

---

## License

FishDB is licensed under the MIT License.

---
