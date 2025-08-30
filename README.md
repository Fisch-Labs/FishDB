# FishDB

FishDB is a lightweight, embeddable, and distributed graph database written in Go, designed for high-performance Retrieval-Augmented Generation (RAG) applications. It aims to be the fastest graph-naive versioned database, offering a powerful way to store, query, and manage evolving graph data.

## Features

*   **Graph-based Data Model:** Store data as nodes and edges, perfect for representing complex relationships and supporting versioning.
*   **High Performance for RAG:** Optimized for rapid data retrieval and augmentation in AI applications.
*   **Naive Versioning:** Efficiently manage and query different versions of your graph data.
*   **Multiple Access Methods:**
    *   RESTful API
    *   GraphQL API (including subscriptions)
*   **Clustering Support:** Distribute your data across multiple nodes for scalability and high availability.
*   **Flexible Storage:**
    *   In-memory storage for high performance.
    *   Disk-based storage for persistence.
*   **Interactive Console:** An interactive command-line console for managing and querying the database.

## Getting Started

### Prerequisites

*   [Go](https://golang.org/doc/install) (version 1.25.0 or later)

### Building from Source

1.  **Clone the repository:**
    ```sh
    git clone https://github.com/Fisch-Labs/FishDB.git
    cd FishDB
    ```

2.  **Tidy up the dependencies:**
    ```sh
    go mod tidy
    ```

3.  **Build the binary:**
    ```sh
    go build ./cli/fishdb.go
    ```

### Running FishDB

You can start the FishDB server using the compiled binary:

```sh
./fishdb -db-path /path/to/your/db
```

This will start the server with a disk-based storage backend at the specified path. For an in-memory database, you can use:

```sh
./fishdb -mem
```

## Usage

### REST API

The REST API provides endpoints for interacting with the database. You can find the available endpoints in the `api/v1` directory.

### GraphQL

FishDB supports GraphQL for querying and mutations. The GraphQL endpoint is typically available at `/graphql`.

### Interactive Console

You can launch the interactive console by running:

```sh
./fishdb -i
```

From the console, you can manage users, and interact with the database directly.

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue.

## License

This project is licensed under the [MIT License](LICENSE).
